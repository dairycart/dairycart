module Products exposing (productsPageHTML, buildProductCards, getProducts)

-- built-ins

import Html exposing (section, node, button, div, h1, text, span, figure, img, p, a, i, header)
import Html.Attributes exposing (class, href, src)
import Html.Events exposing (onClick)
import Http
import Json.Decode exposing (field, string, int, list, Decoder)
import Json.Decode.Pipeline exposing (decode, required, optional)


-- locals

import Models exposing (Product)
import Messages exposing (Msg)


-- dependencies

import Round


-- actual code


fallbackProductImageURL : String
fallbackProductImageURL =
    "https://placehold.it/1280x960"


exampleProducts : List Product
exampleProducts =
    [ buildProduct "Product Name 1" 12.34 123 "SKU1" fallbackProductImageURL
    , buildProduct "Product Name 2" 23.45 234 "SKU2" fallbackProductImageURL
    , buildProduct "Product Name 3" 34.56 345 "SKU3" fallbackProductImageURL
    , buildProduct "Product Name 4" 45.67 456 "SKU4" fallbackProductImageURL
    , buildProduct "Product Name 5" 56.78 567 "SKU5" fallbackProductImageURL
    , buildProduct "Product Name 6" 67.89 678 "SKU6" fallbackProductImageURL
    , buildProduct "Product Name 7" 78.9 789 "SKU7" fallbackProductImageURL
    , buildProduct "Product Name 8" 89.01 890 "SKU8" fallbackProductImageURL
    ]


chunk : Int -> List a -> List (List a)
chunk k xs =
    if List.length xs > k then
        List.take k xs :: chunk k (List.drop k xs)
    else
        [ xs ]


productsPageHTML : List Product -> Html.Html (Cmd Msg)
productsPageHTML products =
    let
        body =
            getProductBody products
    in
        productCardHTML body


getProductBody : List Product -> List (Html.Html (Cmd Msg))
getProductBody products =
    if List.isEmpty products then
        -- [ div []
        --     [ text "no products found. :)"
        --     ]
        -- ]
        [ button [ onClick getProducts ] [ text "load products" ] ]
    else
        -- products
        --     |> chunk 5
        --     |> buildProductCards
        --     |> List.map
        List.map buildProductCards (chunk 5 products)


productCardHTML : List (Html.Html msg) -> Html.Html msg
productCardHTML body =
    div []
        [ section [ class "hero is-small" ]
            [ div [ class "hero-body" ]
                [ div [ class "container" ]
                    [ h1 [ class "title" ]
                        [ text "Manage Products" ]
                    ]
                ]
            ]
        , div [] body
        ]


buildProductCards : List Product -> Html.Html msg
buildProductCards products =
    div [ class "columns" ]
        (List.map buildProductCard products)


buildProductCard : Product -> Html.Html msg
buildProductCard product =
    div [ class "column" ]
        [ div [ class "card" ]
            [ header [ class "card-header" ]
                [ a [ href ("#product/" ++ product.sku) ]
                    [ p [ class "card-header-title" ]
                        [ text product.name ]
                    ]
                ]
            , div [ class "card-image" ]
                [ figure [ class "image is-4by3" ]
                    [ img [ src product.imageURL ]
                        []
                    ]
                ]
            , div [ class "card-content" ]
                [ div [ class "panel-block-item" ]
                    [ span []
                        [ span [ class "icon" ]
                            [ i [ class "fa fa-money" ]
                                []
                            ]
                        , text (" $" ++ Round.round 2 product.price)
                        ]

                    -- , span [ class "is-pulled-right" ]
                    --     [ span [ class "icon" ]
                    --         [ i [ class "fa fa-info" ]
                    --             []
                    --         ]
                    --     , text (toString product.quantity ++ " in stock")
                    --     ]
                    ]
                ]
            ]
        ]


exampleProductURL : String
exampleProductURL =
    "http://api.dairycart.com/v1/product/newborn-sun"


exampleProductsURL : String
exampleProductsURL =
    "http://api.dairycart.com/v1/products"



-- still working on understanding this part


buildProduct : String -> Float -> Int -> String -> String -> Product
buildProduct name price quantity sku image_url =
    { name = name
    , sku = sku
    , price = price
    , quantity = quantity
    , imageURL = image_url
    }


getProducts : Cmd Msg
getProducts =
    list productDecoder
        |> Http.get exampleProductsURL
        |> Http.send Messages.LoadProducts


productDecoder : Decoder Product
productDecoder =
    decode buildProduct
        |> required "name" string
        |> required "price" Json.Decode.float
        |> required "quantity" int
        |> required "sku" string
        |> optional "main_image_url" string fallbackProductImageURL
