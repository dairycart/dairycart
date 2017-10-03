module Dashboard exposing (..)

-- Built-ins

import Html exposing (section, div, h1, h2, text, p, a, button)
import Html.Attributes exposing (class, href, id, style)


-- Local stuff

import Messages exposing (Msg(..))


dashboardHTML : Html.Html Msg
dashboardHTML =
    div []
        [ section [ class "hero is-small" ]
            [ div [ class "hero-body" ]
                [ div [ class "container" ]
                    [ h1 [ class "title" ]
                        [ text "Dairycart Dashboard" ]
                    , h2 [ class "subtitle" ]
                        [ text "Welcome!" ]
                    ]
                ]
            ]
        , section [ class "section" ]
            [ div [ class "columns is-mobile is-multiline" ]
                [ div [ class "column is-half-desktop is-full-mobile" ]
                    [ section [ class "panel" ]
                        [ p [ class "panel-heading" ]
                            [ text "Total Orders" ]
                        , p [ class "panel-tabs" ]
                            [ a [ class "is-active", href "#" ]
                                [ text "Past Week" ]
                            , a [ href "#" ]
                                [ text "Past month" ]
                            , a [ href "#" ]
                                [ text "Past Quarter" ]
                            , a [ href "#" ]
                                [ text "Past Year" ]
                            , a [ href "#" ]
                                [ text "All Time" ]
                            ]
                        , div [ class "panel-block" ]
                            [ div [ id "ordersChart", style [ ( "height", "250px" ) ] ]
                                []
                            ]
                        , div [ class "panel-block" ]
                            [ button [ class "button is-default is-outlined is-fullwidth" ]
                                [ text "View Data" ]
                            ]
                        ]
                    ]
                , div [ class "column is-half-desktop is-full-mobile" ]
                    [ section [ class "panel" ]
                        [ p [ class "panel-heading" ]
                            [ text "Popular Products" ]
                        , div [ class "panel-block" ]
                            [ div [ id "chart2", style [ ( "height", "280px" ) ] ]
                                []
                            ]
                        , div [ class "panel-block" ]
                            [ button [ class "button is-default is-outlined is-fullwidth" ]
                                [ text "View Data" ]
                            ]
                        ]
                    ]
                ]
            ]
        ]
