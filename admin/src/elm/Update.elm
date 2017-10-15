module Update exposing (..)

-- local

import Routing exposing (parseLocation)
import Messages exposing (Msg(..))
import Models exposing (Model)
import Navigation


update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    case msg of
        -- navigation events
        OnLocationChange location ->
            let
                newRoute =
                    parseLocation location
            in
                ( { model | route = newRoute }, Cmd.none )

        GoToProductsPage ->
            ( model, Navigation.newUrl "#products" )

        GoToMainPage ->
            ( model, Navigation.newUrl "#" )

        -- products page events
        LoadProducts (Ok products) ->
            ( { model | products = products }, Cmd.none )

        LoadProducts (Err error) ->
            ( { model | products = [] }, Cmd.none )
