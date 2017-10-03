module Products exposing (productPageHTML, buildProductCards, chunkProducts, products)

-- built-ins

import Html exposing (section, node, button, div, h1, text, span, figure, img, p, a, i, header)
import Html.Attributes exposing (class, href, src)

-- import Http
-- import Json.Decode as Decode


-- dependencies

import Round


products : List Product
products =
    [ { name = "Product Name 1", price = 12.34, quantity = 123, sku = "SKU1", imageURL = "https://placehold.it/1280x960" }
    , { name = "Product Name 2", price = 23.45, quantity = 234, sku = "SKU2", imageURL = "https://placehold.it/1280x960" }
    , { name = "Product Name 3", price = 34.56, quantity = 345, sku = "SKU3", imageURL = "https://placehold.it/1280x960" }
    , { name = "Product Name 4", price = 45.67, quantity = 456, sku = "SKU4", imageURL = "https://placehold.it/1280x960" }
    , { name = "Product Name 5", price = 56.78, quantity = 567, sku = "SKU5", imageURL = "https://placehold.it/1280x960" }
    , { name = "Product Name 6", price = 67.89, quantity = 678, sku = "SKU6", imageURL = "https://placehold.it/1280x960" }
    , { name = "Product Name 7", price = 78.90, quantity = 789, sku = "SKU7", imageURL = "https://placehold.it/1280x960" }
    , { name = "Product Name 8", price = 89.01, quantity = 890, sku = "SKU8", imageURL = "https://placehold.it/1280x960" }
    ]


type alias Product =
    { name : String
    , sku : String
    , price : Float
    , quantity : Int
    , imageURL : String
    }


chunkProducts : Int -> List a -> List (List a)
chunkProducts k xs =
    if List.length xs > k then
        List.take k xs :: chunkProducts k (List.drop k xs)
    else
        [ xs ]


productPageHTML : List (Html.Html msg) -> Html.Html msg
productPageHTML productCards =
    div []
        [ section [ class "hero is-small" ]
            [ div [ class "hero-body" ]
                [ div [ class "container" ]
                    [ h1 [ class "title" ]
                        [ text "Manage Products" ]
                    ]
                ]
            ]
        , div [] productCards
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
                [ a [ href ("product/" ++ product.sku) ]
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


-- still working on understanding this part

-- getProduct : Product
-- getProduct =
--     let
--         req = Http.get "http://localhost:4321/v1/product/newborn-sun"
--     in
--         decodeProduct req
-- getProduct : List(Product)
-- getProduct =
--     let
--         req = Http.get "http://localhost:4321/v1/product/newborn-sun"
--     in
--         [(decodeProduct req)]
-- decodeProduct : Decode.Decoder Product
-- decodeProduct req =
--     Decode.map2 Product
--         (Decode.field "name" Decode.string)
--         (Decode.field "price" Decode.float)
