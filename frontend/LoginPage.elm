module LoginPage
    exposing
        ( Model
        , Msg
        , Context
        , init
        , update
        , view
        )

import HttpBuilder as Http
import Task exposing (Task)
import Json.Encode
import Json.Decode exposing ((:=), Decoder)
import Html
    exposing
        ( Html
        , div
        , input
        , button
        , text
        , form
        )
import Html.Events
    exposing
        ( onInput
        , onClick
        , onSubmit
        )
import Html.Attributes
    exposing
        ( placeholder
        , value
        , class
        , type'
        , disabled
        )


-- MODEL


type alias Model =
    { errorResponse : String
    , loginForm : LoginForm
    }


init : Model
init =
    { errorResponse = ""
    , loginForm = emptyForm
    }


type alias LoginForm =
    { username : String
    , password : String
    }


emptyForm : LoginForm
emptyForm =
    LoginForm "" ""


type alias Context msg =
    {-
       url configures the URL of the login backend

       onSuccess accepts a message which will be run
       if the login succeeds. For example RedirectTo url.
    -}
    { url : String
    , onSuccess : Maybe (Cmd msg)
    }



-- UPDATE


type Msg
    = NoOp
    | SetUsername String
    | SetPassword String
    | SubmitLogin
    | SubmitFail (Http.Error String)
    | SubmitSucceed (Http.Response String)


update : Context Msg -> Msg -> Model -> ( Model, Cmd Msg )
update ctx msg model =
    case msg of
        NoOp ->
            model ! []

        SetUsername newInput ->
            let
                setUsername user =
                    { user | username = newInput }
            in
                { model | loginForm = setUsername model.loginForm } ! []

        SetPassword newInput ->
            let
                setPassword form =
                    { form | password = newInput }
            in
                { model | loginForm = setPassword model.loginForm } ! []

        SubmitLogin ->
            let
                loginForm : LoginForm
                loginForm =
                    model.loginForm

                encodeUser : Json.Encode.Value
                encodeUser =
                    Json.Encode.object
                        [ ( "username", Json.Encode.string loginForm.username )
                        , ( "password", Json.Encode.string loginForm.password )
                        ]

                loginUser : Task (Http.Error String) (Http.Response String)
                loginUser =
                    Http.post ctx.url
                        |> Http.withHeader "Content-Type" "application/json"
                        |> Http.withJsonBody encodeUser
                        |> Http.send (Http.stringReader) Http.stringReader

                login : Cmd Msg
                login =
                    Task.perform SubmitFail SubmitSucceed loginUser
            in
                model ! [ login ]

        SubmitFail err ->
            let
                response =
                    case err of
                        Http.BadResponse resp ->
                            resp.data

                        _ ->
                            -- TODO
                            "unknown error"
            in
                { model | errorResponse = response } ! []

        SubmitSucceed response ->
            case ctx.onSuccess of
                Nothing ->
                    { model | errorResponse = "Success" } ! []

                Just cmdMsg ->
                    { model | errorResponse = "Success" } ! [ cmdMsg ]



-- VIEW


view : Model -> Html Msg
view model =
    let
        errorMessage : Html Msg
        errorMessage =
            if model.errorResponse /= "" then
                div [ class "loginPage__form__error" ]
                    [ text model.errorResponse ]
            else
                emptyNode

        disableSubmit : Bool
        disableSubmit =
            model.loginForm.username == "" || model.loginForm.password == ""

        usernameField =
            input
                [ placeholder "username"
                , value model.loginForm.username
                , onInput SetUsername
                , class "loginPage__form__username"
                ]
                []

        passwordField =
            input
                [ placeholder "password"
                , value model.loginForm.password
                , onInput SetPassword
                , class "loginPage__form__password"
                , type' "password"
                ]
                []

        submitButton =
            button
                [ class "loginPage__form__button"
                , disabled disableSubmit
                ]
                [ text "LOGIN" ]

        loginForm =
            form
                [ class "loginPage__form"
                , onSubmit SubmitLogin
                ]
                [ usernameField
                , errorMessage
                , passwordField
                , submitButton
                ]
    in
        div
            [ class "loginPage" ]
            [ loginForm ]


emptyNode : Html Msg
emptyNode =
    text ""
