module Messages exposing (..)

import Navigation exposing (Location)
import Http


-- locals

import Models exposing (Product)


type Msg
    = OnLocationChange Location
    | GoToProductsPage
    | GoToMainPage
    | LoadProducts (Result Http.Error (List Product))
