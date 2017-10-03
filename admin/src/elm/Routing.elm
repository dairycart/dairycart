module Routing exposing (..)

import Navigation exposing (Location)
import UrlParser exposing (..)


type Route
    = MainPage
    | ProductsPage
    | NotFoundRoute


matchers : Parser (Route -> a) a
matchers =
    oneOf
        [ map MainPage top
        , map ProductsPage (s "products")
        ]


parseLocation : Location -> Route
parseLocation location =
    case (parseHash matchers location) of
        Just route ->
            route

        Nothing ->
            NotFoundRoute
