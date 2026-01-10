package middleware

import (
	"context"
	"net/http"

	"golang.org/x/text/language"

	"github.com/jacobpq/soccer-manager/internal/locales"
)

func Locale(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		matcher := locales.GetMatcher()

		accept := r.Header.Get("Accept-Language")
		tag, _ := language.MatchStrings(matcher, accept)

		base, _ := tag.Base()
		lang := base.String()

		ctx := context.WithValue(r.Context(), locales.Key, lang)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
