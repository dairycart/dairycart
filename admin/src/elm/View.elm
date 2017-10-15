module View exposing (..)

import Html exposing (Html, div, h1, text)
import Messages exposing (Msg(..))
import Models exposing (Model, Product)
import Routing exposing (Route(..))


-- Local Stuff

import Products exposing (getProducts, productsPageHTML)
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

        ProductPage sku ->
            notFoundView

        ProductsPage ->
            productsPage model

        NotFoundRoute ->
            notFoundView


mainPage : Html Msg
mainPage =
    dashboardHTML


productsPage : { a | products : List Product } -> Html Msg
productsPage model =
    productsPageHTML model.products


notFoundView : Html Msg
notFoundView =
    h1 []
        [ text "Not Found" ]
