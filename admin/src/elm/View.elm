module View exposing (..)

import Html exposing (Html, div, text)
import Messages exposing (Msg(..))
import Models exposing (Model)
import Routing exposing (Route(..))


-- Local Stuff

import Products exposing (productPageHTML, buildProductCards, chunkProducts, products)
import Dashboard exposing (dashboardHTML)


view : Model -> Html Msg
view model =
    div []
        [ page model ]


page : Model -> Html Msg
page model =
    case model.route of
        MainPage ->
            mainPage

        ProductsPage ->
            productsPage

        NotFoundRoute ->
            notFoundView


mainPage : Html Msg
mainPage =
    dashboardHTML


productsPage : Html Msg
productsPage =
    productPageHTML (List.map buildProductCards (chunkProducts 5 products))


notFoundView : Html Msg
notFoundView =
    div []
        [ text "Not Found" ]
