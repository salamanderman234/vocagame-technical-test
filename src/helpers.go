package src

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/asaskevich/govalidator"
	"github.com/golang-jwt/jwt/v5"
	"github.com/microcosm-cc/bluemonday"
)

// response
type Response struct {
	Msg     string `json:"message"`
	Details any    `json:"details,omitempty"`
}

func SendJSON(w http.ResponseWriter, status int, msg string, details any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	data := Response{
		Msg:     msg,
		Details: details,
	}
	json.NewEncoder(w).Encode(data)
}

// end of response
// validation
func ValidationErrorsExtration(err error) any {
	cont := []any{}
	errs, ok := err.(govalidator.Errors)
	if !ok {
		return nil
	}
	for _, e := range errs {
		eCon, ok := e.(govalidator.Error)
		if !ok {
			continue
		}
		detail := struct {
			Field   string `json:"field"`
			Rule    string `json:"rule"`
			Message string `json:"message"`
		}{
			Field:   eCon.Name,
			Rule:    eCon.Validator,
			Message: eCon.Error(),
		}
		cont = append(cont, detail)
	}
	return cont
}

func ValidateStruct(data any) error {
	result, err := govalidator.ValidateStruct(data)
	if !result {
		copy := ErrValidation
		copy.Obj = err
		return copy
	}
	return nil
}

// end of validation
// policy
type Key string

const UserContextKey Key = "user"

func GetUserContext(ctx context.Context) (User, error) {
	val := ctx.Value(UserContextKey)
	user, ok := val.(User)
	if !ok && val == "" {
		return user, ErrAuthorize
	}
	return user, nil
}

type AuthorizeFunc func(user User, data any) bool

func Authorize(ctx context.Context, data any, fun AuthorizeFunc) error {
	user, err := GetUserContext(ctx)
	if err != nil {
		return err
	}
	result := fun(user, data)
	if !result {
		return ErrAuthorize
	}
	return nil
}

// end of policy
// request
func Sanitize(text string) string {
	policy := bluemonday.UGCPolicy()
	return policy.Sanitize(text)
}
func GetRequestContext(r *http.Request) context.Context {
	ctx := r.Context()
	return ctx
}

func SetRequestContext(r *http.Request, key Key, val any) *http.Request {
	ctx := r.Context()
	return r.WithContext(context.WithValue(ctx, key, val))
}

func ErrorHandler(w http.ResponseWriter, err error) {
	errCon, ok := err.(Err)
	errs, errsOk := errCon.Obj.(govalidator.Errors)
	if !ok {
		SendJSON(w, http.StatusInternalServerError, err.Error(), nil)
	} else {
		var details any
		status := errCon.Status
		message := errCon.Error()
		if errsOk {
			details = ValidationErrorsExtration(errs)
		}
		SendJSON(w, status, message, details)
	}

}

// end of request
// encrypt
func EncryptData(text string) (string, error) {
	c := CiperBlock
	cipherText := make([]byte, aes.BlockSize+len(text))
	iv := cipherText[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}
	stream := cipher.NewCFBEncrypter(c, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], []byte(text))
	return base64.URLEncoding.EncodeToString(cipherText), nil
}

func DecryptData(cipherText string) (string, error) {
	decodedCiphertext, err := base64.URLEncoding.DecodeString(cipherText)
	if err != nil {
		return "", err
	}
	block := CiperBlock
	iv := decodedCiphertext[:aes.BlockSize]
	decodedCiphertext = decodedCiphertext[aes.BlockSize:]
	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(decodedCiphertext, decodedCiphertext)
	return string(decodedCiphertext), nil
}

// end of encrypt
// jwt
func CreateToken(user User) (string, error) {
	idStr := strconv.Itoa(int(user.ID))
	claims := JWTClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ID:      idStr,
			Subject: "Auth token",
		},
		Email: user.Email,
		Role:  user.Role,
	}
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)
	return token.SignedString([]byte(JWT_KEY))
}
func ParseToken(token string) (User, error) {
	user := User{}
	tkn, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		method := t.Method
		if method != jwt.SigningMethodHS256 {
			return nil, jwt.ErrTokenSignatureInvalid
		}
		return []byte(JWT_KEY), nil
	})
	if err != nil {
		return user, err
	}
	if !tkn.Valid {
		return user, ErrAuthorize
	}
	var result JWTClaims
	claims, _ := tkn.Claims.(jwt.MapClaims)
	err = TranslateStruct(claims, &result)
	if err != nil {
		return user, err
	}
	id, _ := strconv.Atoi(result.ID)
	user.ID = uint(id)
	user.Email = result.Email
	user.Role = result.Role
	return user, nil
}

// end of jwt
// translate
func TranslateStruct(source any, target any) error {
	jsonEncode, err := json.Marshal(source)
	if err != nil {
		return err
	}
	err = json.Unmarshal(jsonEncode, target)
	if err != nil {
		return err
	}
	return nil
}

// end of translate
