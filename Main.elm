module Main exposing (Model, init, update, view, subscriptions)

import Html.App as App
import Html exposing (Html, div)
import LoginPage as Login


-- MAIN


main : Program Never
main =
    App.program
        { init = ( init, Cmd.none )
        , update = update
        , view = view
        , subscriptions = subscriptions
        }


subscriptions : Model -> Sub Msg
subscriptions model =
    Sub.none



-- MODEL


type alias Model =
    { login : Login.Model
    }


init : Model
init =
    { login = Login.init
    }



-- UPDATE


type Msg
    = NoOp
    | LoginMsg Login.Msg


update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    case Debug.log "msg:" msg of
        NoOp ->
            model ! []

        LoginMsg loginMsg ->
            let
                loginCtx : Login.Context msg
                loginCtx =
                    { url = "http://localhost:8000/login"
                    , onSuccess = Nothing
                    }

                ( newLogin, newMsg ) =
                    Login.update loginCtx loginMsg model.login
            in
                { model | login = newLogin } ! [ Cmd.map LoginMsg newMsg ]



-- VIEW


view : Model -> Html Msg
view model =
    let
        login =
            App.map LoginMsg <| Login.view model.login
    in
        div [] [ login ]
