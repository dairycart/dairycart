module Models exposing (Model, Product, initialModel)

import Routing


type alias Model =
    { route : Routing.Route
    , products : List Product
    }


type alias Product =
    { name : String
    , sku : String
    , price : Float
    , quantity : Int
    , imageURL : String
    }


initialModel : Routing.Route -> Model
initialModel route =
    { route = route, products = [] }
