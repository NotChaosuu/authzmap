package main

// EXPECT: 400 modules via an in-tree Router interface + depth-40 helper chain; open module m399 DELETE /m399/items/{id} flagged; complete or INCOMPLETE (never silently truncated); < 30 s
import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

func valid(s string) bool { return s == "x" }
func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !valid(r.Header.Get("Authorization")) {
			http.Error(w, "unauthorized", 401)
			return
		}
		next.ServeHTTP(w, r)
	})
}
func h(w http.ResponseWriter, r *http.Request) {}

type Router interface {
	Get(pattern string, fn http.HandlerFunc)
	Post(pattern string, fn http.HandlerFunc)
	Delete(pattern string, fn http.HandlerFunc)
	With(m ...func(http.Handler) http.Handler) chi.Router
	Route(p string, fn func(chi.Router)) chi.Router
}
type Module interface{ Register(r Router) }
type m0 struct{}

func (m0) Register(r Router) { reg0(r, "/m0") }
func reg0(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain0_0(r, p)
}
func chain0_0(r Router, p string) { r.Get(p+"/deep0", h); chain0_1(r, p) }
func chain0_1(r Router, p string) { r.Get(p+"/deep1", h); chain0_2(r, p) }
func chain0_2(r Router, p string) { r.Get(p+"/deep2", h); chain0_3(r, p) }
func chain0_3(r Router, p string) { deep_0(r, p) }

type m1 struct{}

func (m1) Register(r Router) { reg1(r, "/m1") }
func reg1(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain1_0(r, p)
}
func chain1_0(r Router, p string) { r.Get(p+"/deep0", h); chain1_1(r, p) }
func chain1_1(r Router, p string) { r.Get(p+"/deep1", h); chain1_2(r, p) }
func chain1_2(r Router, p string) { r.Get(p+"/deep2", h); chain1_3(r, p) }
func chain1_3(r Router, p string) { deep_0(r, p) }

type m2 struct{}

func (m2) Register(r Router) { reg2(r, "/m2") }
func reg2(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain2_0(r, p)
}
func chain2_0(r Router, p string) { r.Get(p+"/deep0", h); chain2_1(r, p) }
func chain2_1(r Router, p string) { r.Get(p+"/deep1", h); chain2_2(r, p) }
func chain2_2(r Router, p string) { r.Get(p+"/deep2", h); chain2_3(r, p) }
func chain2_3(r Router, p string) { deep_0(r, p) }

type m3 struct{}

func (m3) Register(r Router) { reg3(r, "/m3") }
func reg3(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain3_0(r, p)
}
func chain3_0(r Router, p string) { r.Get(p+"/deep0", h); chain3_1(r, p) }
func chain3_1(r Router, p string) { r.Get(p+"/deep1", h); chain3_2(r, p) }
func chain3_2(r Router, p string) { r.Get(p+"/deep2", h); chain3_3(r, p) }
func chain3_3(r Router, p string) { deep_0(r, p) }

type m4 struct{}

func (m4) Register(r Router) { reg4(r, "/m4") }
func reg4(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain4_0(r, p)
}
func chain4_0(r Router, p string) { r.Get(p+"/deep0", h); chain4_1(r, p) }
func chain4_1(r Router, p string) { r.Get(p+"/deep1", h); chain4_2(r, p) }
func chain4_2(r Router, p string) { r.Get(p+"/deep2", h); chain4_3(r, p) }
func chain4_3(r Router, p string) { deep_0(r, p) }

type m5 struct{}

func (m5) Register(r Router) { reg5(r, "/m5") }
func reg5(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain5_0(r, p)
}
func chain5_0(r Router, p string) { r.Get(p+"/deep0", h); chain5_1(r, p) }
func chain5_1(r Router, p string) { r.Get(p+"/deep1", h); chain5_2(r, p) }
func chain5_2(r Router, p string) { r.Get(p+"/deep2", h); chain5_3(r, p) }
func chain5_3(r Router, p string) { deep_0(r, p) }

type m6 struct{}

func (m6) Register(r Router) { reg6(r, "/m6") }
func reg6(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain6_0(r, p)
}
func chain6_0(r Router, p string) { r.Get(p+"/deep0", h); chain6_1(r, p) }
func chain6_1(r Router, p string) { r.Get(p+"/deep1", h); chain6_2(r, p) }
func chain6_2(r Router, p string) { r.Get(p+"/deep2", h); chain6_3(r, p) }
func chain6_3(r Router, p string) { deep_0(r, p) }

type m7 struct{}

func (m7) Register(r Router) { reg7(r, "/m7") }
func reg7(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain7_0(r, p)
}
func chain7_0(r Router, p string) { r.Get(p+"/deep0", h); chain7_1(r, p) }
func chain7_1(r Router, p string) { r.Get(p+"/deep1", h); chain7_2(r, p) }
func chain7_2(r Router, p string) { r.Get(p+"/deep2", h); chain7_3(r, p) }
func chain7_3(r Router, p string) { deep_0(r, p) }

type m8 struct{}

func (m8) Register(r Router) { reg8(r, "/m8") }
func reg8(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain8_0(r, p)
}
func chain8_0(r Router, p string) { r.Get(p+"/deep0", h); chain8_1(r, p) }
func chain8_1(r Router, p string) { r.Get(p+"/deep1", h); chain8_2(r, p) }
func chain8_2(r Router, p string) { r.Get(p+"/deep2", h); chain8_3(r, p) }
func chain8_3(r Router, p string) { deep_0(r, p) }

type m9 struct{}

func (m9) Register(r Router) { reg9(r, "/m9") }
func reg9(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain9_0(r, p)
}
func chain9_0(r Router, p string) { r.Get(p+"/deep0", h); chain9_1(r, p) }
func chain9_1(r Router, p string) { r.Get(p+"/deep1", h); chain9_2(r, p) }
func chain9_2(r Router, p string) { r.Get(p+"/deep2", h); chain9_3(r, p) }
func chain9_3(r Router, p string) { deep_0(r, p) }

type m10 struct{}

func (m10) Register(r Router) { reg10(r, "/m10") }
func reg10(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain10_0(r, p)
}
func chain10_0(r Router, p string) { r.Get(p+"/deep0", h); chain10_1(r, p) }
func chain10_1(r Router, p string) { r.Get(p+"/deep1", h); chain10_2(r, p) }
func chain10_2(r Router, p string) { r.Get(p+"/deep2", h); chain10_3(r, p) }
func chain10_3(r Router, p string) { deep_0(r, p) }

type m11 struct{}

func (m11) Register(r Router) { reg11(r, "/m11") }
func reg11(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain11_0(r, p)
}
func chain11_0(r Router, p string) { r.Get(p+"/deep0", h); chain11_1(r, p) }
func chain11_1(r Router, p string) { r.Get(p+"/deep1", h); chain11_2(r, p) }
func chain11_2(r Router, p string) { r.Get(p+"/deep2", h); chain11_3(r, p) }
func chain11_3(r Router, p string) { deep_0(r, p) }

type m12 struct{}

func (m12) Register(r Router) { reg12(r, "/m12") }
func reg12(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain12_0(r, p)
}
func chain12_0(r Router, p string) { r.Get(p+"/deep0", h); chain12_1(r, p) }
func chain12_1(r Router, p string) { r.Get(p+"/deep1", h); chain12_2(r, p) }
func chain12_2(r Router, p string) { r.Get(p+"/deep2", h); chain12_3(r, p) }
func chain12_3(r Router, p string) { deep_0(r, p) }

type m13 struct{}

func (m13) Register(r Router) { reg13(r, "/m13") }
func reg13(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain13_0(r, p)
}
func chain13_0(r Router, p string) { r.Get(p+"/deep0", h); chain13_1(r, p) }
func chain13_1(r Router, p string) { r.Get(p+"/deep1", h); chain13_2(r, p) }
func chain13_2(r Router, p string) { r.Get(p+"/deep2", h); chain13_3(r, p) }
func chain13_3(r Router, p string) { deep_0(r, p) }

type m14 struct{}

func (m14) Register(r Router) { reg14(r, "/m14") }
func reg14(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain14_0(r, p)
}
func chain14_0(r Router, p string) { r.Get(p+"/deep0", h); chain14_1(r, p) }
func chain14_1(r Router, p string) { r.Get(p+"/deep1", h); chain14_2(r, p) }
func chain14_2(r Router, p string) { r.Get(p+"/deep2", h); chain14_3(r, p) }
func chain14_3(r Router, p string) { deep_0(r, p) }

type m15 struct{}

func (m15) Register(r Router) { reg15(r, "/m15") }
func reg15(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain15_0(r, p)
}
func chain15_0(r Router, p string) { r.Get(p+"/deep0", h); chain15_1(r, p) }
func chain15_1(r Router, p string) { r.Get(p+"/deep1", h); chain15_2(r, p) }
func chain15_2(r Router, p string) { r.Get(p+"/deep2", h); chain15_3(r, p) }
func chain15_3(r Router, p string) { deep_0(r, p) }

type m16 struct{}

func (m16) Register(r Router) { reg16(r, "/m16") }
func reg16(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain16_0(r, p)
}
func chain16_0(r Router, p string) { r.Get(p+"/deep0", h); chain16_1(r, p) }
func chain16_1(r Router, p string) { r.Get(p+"/deep1", h); chain16_2(r, p) }
func chain16_2(r Router, p string) { r.Get(p+"/deep2", h); chain16_3(r, p) }
func chain16_3(r Router, p string) { deep_0(r, p) }

type m17 struct{}

func (m17) Register(r Router) { reg17(r, "/m17") }
func reg17(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain17_0(r, p)
}
func chain17_0(r Router, p string) { r.Get(p+"/deep0", h); chain17_1(r, p) }
func chain17_1(r Router, p string) { r.Get(p+"/deep1", h); chain17_2(r, p) }
func chain17_2(r Router, p string) { r.Get(p+"/deep2", h); chain17_3(r, p) }
func chain17_3(r Router, p string) { deep_0(r, p) }

type m18 struct{}

func (m18) Register(r Router) { reg18(r, "/m18") }
func reg18(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain18_0(r, p)
}
func chain18_0(r Router, p string) { r.Get(p+"/deep0", h); chain18_1(r, p) }
func chain18_1(r Router, p string) { r.Get(p+"/deep1", h); chain18_2(r, p) }
func chain18_2(r Router, p string) { r.Get(p+"/deep2", h); chain18_3(r, p) }
func chain18_3(r Router, p string) { deep_0(r, p) }

type m19 struct{}

func (m19) Register(r Router) { reg19(r, "/m19") }
func reg19(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain19_0(r, p)
}
func chain19_0(r Router, p string) { r.Get(p+"/deep0", h); chain19_1(r, p) }
func chain19_1(r Router, p string) { r.Get(p+"/deep1", h); chain19_2(r, p) }
func chain19_2(r Router, p string) { r.Get(p+"/deep2", h); chain19_3(r, p) }
func chain19_3(r Router, p string) { deep_0(r, p) }

type m20 struct{}

func (m20) Register(r Router) { reg20(r, "/m20") }
func reg20(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain20_0(r, p)
}
func chain20_0(r Router, p string) { r.Get(p+"/deep0", h); chain20_1(r, p) }
func chain20_1(r Router, p string) { r.Get(p+"/deep1", h); chain20_2(r, p) }
func chain20_2(r Router, p string) { r.Get(p+"/deep2", h); chain20_3(r, p) }
func chain20_3(r Router, p string) { deep_0(r, p) }

type m21 struct{}

func (m21) Register(r Router) { reg21(r, "/m21") }
func reg21(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain21_0(r, p)
}
func chain21_0(r Router, p string) { r.Get(p+"/deep0", h); chain21_1(r, p) }
func chain21_1(r Router, p string) { r.Get(p+"/deep1", h); chain21_2(r, p) }
func chain21_2(r Router, p string) { r.Get(p+"/deep2", h); chain21_3(r, p) }
func chain21_3(r Router, p string) { deep_0(r, p) }

type m22 struct{}

func (m22) Register(r Router) { reg22(r, "/m22") }
func reg22(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain22_0(r, p)
}
func chain22_0(r Router, p string) { r.Get(p+"/deep0", h); chain22_1(r, p) }
func chain22_1(r Router, p string) { r.Get(p+"/deep1", h); chain22_2(r, p) }
func chain22_2(r Router, p string) { r.Get(p+"/deep2", h); chain22_3(r, p) }
func chain22_3(r Router, p string) { deep_0(r, p) }

type m23 struct{}

func (m23) Register(r Router) { reg23(r, "/m23") }
func reg23(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain23_0(r, p)
}
func chain23_0(r Router, p string) { r.Get(p+"/deep0", h); chain23_1(r, p) }
func chain23_1(r Router, p string) { r.Get(p+"/deep1", h); chain23_2(r, p) }
func chain23_2(r Router, p string) { r.Get(p+"/deep2", h); chain23_3(r, p) }
func chain23_3(r Router, p string) { deep_0(r, p) }

type m24 struct{}

func (m24) Register(r Router) { reg24(r, "/m24") }
func reg24(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain24_0(r, p)
}
func chain24_0(r Router, p string) { r.Get(p+"/deep0", h); chain24_1(r, p) }
func chain24_1(r Router, p string) { r.Get(p+"/deep1", h); chain24_2(r, p) }
func chain24_2(r Router, p string) { r.Get(p+"/deep2", h); chain24_3(r, p) }
func chain24_3(r Router, p string) { deep_0(r, p) }

type m25 struct{}

func (m25) Register(r Router) { reg25(r, "/m25") }
func reg25(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain25_0(r, p)
}
func chain25_0(r Router, p string) { r.Get(p+"/deep0", h); chain25_1(r, p) }
func chain25_1(r Router, p string) { r.Get(p+"/deep1", h); chain25_2(r, p) }
func chain25_2(r Router, p string) { r.Get(p+"/deep2", h); chain25_3(r, p) }
func chain25_3(r Router, p string) { deep_0(r, p) }

type m26 struct{}

func (m26) Register(r Router) { reg26(r, "/m26") }
func reg26(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain26_0(r, p)
}
func chain26_0(r Router, p string) { r.Get(p+"/deep0", h); chain26_1(r, p) }
func chain26_1(r Router, p string) { r.Get(p+"/deep1", h); chain26_2(r, p) }
func chain26_2(r Router, p string) { r.Get(p+"/deep2", h); chain26_3(r, p) }
func chain26_3(r Router, p string) { deep_0(r, p) }

type m27 struct{}

func (m27) Register(r Router) { reg27(r, "/m27") }
func reg27(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain27_0(r, p)
}
func chain27_0(r Router, p string) { r.Get(p+"/deep0", h); chain27_1(r, p) }
func chain27_1(r Router, p string) { r.Get(p+"/deep1", h); chain27_2(r, p) }
func chain27_2(r Router, p string) { r.Get(p+"/deep2", h); chain27_3(r, p) }
func chain27_3(r Router, p string) { deep_0(r, p) }

type m28 struct{}

func (m28) Register(r Router) { reg28(r, "/m28") }
func reg28(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain28_0(r, p)
}
func chain28_0(r Router, p string) { r.Get(p+"/deep0", h); chain28_1(r, p) }
func chain28_1(r Router, p string) { r.Get(p+"/deep1", h); chain28_2(r, p) }
func chain28_2(r Router, p string) { r.Get(p+"/deep2", h); chain28_3(r, p) }
func chain28_3(r Router, p string) { deep_0(r, p) }

type m29 struct{}

func (m29) Register(r Router) { reg29(r, "/m29") }
func reg29(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain29_0(r, p)
}
func chain29_0(r Router, p string) { r.Get(p+"/deep0", h); chain29_1(r, p) }
func chain29_1(r Router, p string) { r.Get(p+"/deep1", h); chain29_2(r, p) }
func chain29_2(r Router, p string) { r.Get(p+"/deep2", h); chain29_3(r, p) }
func chain29_3(r Router, p string) { deep_0(r, p) }

type m30 struct{}

func (m30) Register(r Router) { reg30(r, "/m30") }
func reg30(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain30_0(r, p)
}
func chain30_0(r Router, p string) { r.Get(p+"/deep0", h); chain30_1(r, p) }
func chain30_1(r Router, p string) { r.Get(p+"/deep1", h); chain30_2(r, p) }
func chain30_2(r Router, p string) { r.Get(p+"/deep2", h); chain30_3(r, p) }
func chain30_3(r Router, p string) { deep_0(r, p) }

type m31 struct{}

func (m31) Register(r Router) { reg31(r, "/m31") }
func reg31(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain31_0(r, p)
}
func chain31_0(r Router, p string) { r.Get(p+"/deep0", h); chain31_1(r, p) }
func chain31_1(r Router, p string) { r.Get(p+"/deep1", h); chain31_2(r, p) }
func chain31_2(r Router, p string) { r.Get(p+"/deep2", h); chain31_3(r, p) }
func chain31_3(r Router, p string) { deep_0(r, p) }

type m32 struct{}

func (m32) Register(r Router) { reg32(r, "/m32") }
func reg32(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain32_0(r, p)
}
func chain32_0(r Router, p string) { r.Get(p+"/deep0", h); chain32_1(r, p) }
func chain32_1(r Router, p string) { r.Get(p+"/deep1", h); chain32_2(r, p) }
func chain32_2(r Router, p string) { r.Get(p+"/deep2", h); chain32_3(r, p) }
func chain32_3(r Router, p string) { deep_0(r, p) }

type m33 struct{}

func (m33) Register(r Router) { reg33(r, "/m33") }
func reg33(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain33_0(r, p)
}
func chain33_0(r Router, p string) { r.Get(p+"/deep0", h); chain33_1(r, p) }
func chain33_1(r Router, p string) { r.Get(p+"/deep1", h); chain33_2(r, p) }
func chain33_2(r Router, p string) { r.Get(p+"/deep2", h); chain33_3(r, p) }
func chain33_3(r Router, p string) { deep_0(r, p) }

type m34 struct{}

func (m34) Register(r Router) { reg34(r, "/m34") }
func reg34(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain34_0(r, p)
}
func chain34_0(r Router, p string) { r.Get(p+"/deep0", h); chain34_1(r, p) }
func chain34_1(r Router, p string) { r.Get(p+"/deep1", h); chain34_2(r, p) }
func chain34_2(r Router, p string) { r.Get(p+"/deep2", h); chain34_3(r, p) }
func chain34_3(r Router, p string) { deep_0(r, p) }

type m35 struct{}

func (m35) Register(r Router) { reg35(r, "/m35") }
func reg35(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain35_0(r, p)
}
func chain35_0(r Router, p string) { r.Get(p+"/deep0", h); chain35_1(r, p) }
func chain35_1(r Router, p string) { r.Get(p+"/deep1", h); chain35_2(r, p) }
func chain35_2(r Router, p string) { r.Get(p+"/deep2", h); chain35_3(r, p) }
func chain35_3(r Router, p string) { deep_0(r, p) }

type m36 struct{}

func (m36) Register(r Router) { reg36(r, "/m36") }
func reg36(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain36_0(r, p)
}
func chain36_0(r Router, p string) { r.Get(p+"/deep0", h); chain36_1(r, p) }
func chain36_1(r Router, p string) { r.Get(p+"/deep1", h); chain36_2(r, p) }
func chain36_2(r Router, p string) { r.Get(p+"/deep2", h); chain36_3(r, p) }
func chain36_3(r Router, p string) { deep_0(r, p) }

type m37 struct{}

func (m37) Register(r Router) { reg37(r, "/m37") }
func reg37(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain37_0(r, p)
}
func chain37_0(r Router, p string) { r.Get(p+"/deep0", h); chain37_1(r, p) }
func chain37_1(r Router, p string) { r.Get(p+"/deep1", h); chain37_2(r, p) }
func chain37_2(r Router, p string) { r.Get(p+"/deep2", h); chain37_3(r, p) }
func chain37_3(r Router, p string) { deep_0(r, p) }

type m38 struct{}

func (m38) Register(r Router) { reg38(r, "/m38") }
func reg38(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain38_0(r, p)
}
func chain38_0(r Router, p string) { r.Get(p+"/deep0", h); chain38_1(r, p) }
func chain38_1(r Router, p string) { r.Get(p+"/deep1", h); chain38_2(r, p) }
func chain38_2(r Router, p string) { r.Get(p+"/deep2", h); chain38_3(r, p) }
func chain38_3(r Router, p string) { deep_0(r, p) }

type m39 struct{}

func (m39) Register(r Router) { reg39(r, "/m39") }
func reg39(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain39_0(r, p)
}
func chain39_0(r Router, p string) { r.Get(p+"/deep0", h); chain39_1(r, p) }
func chain39_1(r Router, p string) { r.Get(p+"/deep1", h); chain39_2(r, p) }
func chain39_2(r Router, p string) { r.Get(p+"/deep2", h); chain39_3(r, p) }
func chain39_3(r Router, p string) { deep_0(r, p) }

type m40 struct{}

func (m40) Register(r Router) { reg40(r, "/m40") }
func reg40(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain40_0(r, p)
}
func chain40_0(r Router, p string) { r.Get(p+"/deep0", h); chain40_1(r, p) }
func chain40_1(r Router, p string) { r.Get(p+"/deep1", h); chain40_2(r, p) }
func chain40_2(r Router, p string) { r.Get(p+"/deep2", h); chain40_3(r, p) }
func chain40_3(r Router, p string) { deep_0(r, p) }

type m41 struct{}

func (m41) Register(r Router) { reg41(r, "/m41") }
func reg41(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain41_0(r, p)
}
func chain41_0(r Router, p string) { r.Get(p+"/deep0", h); chain41_1(r, p) }
func chain41_1(r Router, p string) { r.Get(p+"/deep1", h); chain41_2(r, p) }
func chain41_2(r Router, p string) { r.Get(p+"/deep2", h); chain41_3(r, p) }
func chain41_3(r Router, p string) { deep_0(r, p) }

type m42 struct{}

func (m42) Register(r Router) { reg42(r, "/m42") }
func reg42(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain42_0(r, p)
}
func chain42_0(r Router, p string) { r.Get(p+"/deep0", h); chain42_1(r, p) }
func chain42_1(r Router, p string) { r.Get(p+"/deep1", h); chain42_2(r, p) }
func chain42_2(r Router, p string) { r.Get(p+"/deep2", h); chain42_3(r, p) }
func chain42_3(r Router, p string) { deep_0(r, p) }

type m43 struct{}

func (m43) Register(r Router) { reg43(r, "/m43") }
func reg43(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain43_0(r, p)
}
func chain43_0(r Router, p string) { r.Get(p+"/deep0", h); chain43_1(r, p) }
func chain43_1(r Router, p string) { r.Get(p+"/deep1", h); chain43_2(r, p) }
func chain43_2(r Router, p string) { r.Get(p+"/deep2", h); chain43_3(r, p) }
func chain43_3(r Router, p string) { deep_0(r, p) }

type m44 struct{}

func (m44) Register(r Router) { reg44(r, "/m44") }
func reg44(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain44_0(r, p)
}
func chain44_0(r Router, p string) { r.Get(p+"/deep0", h); chain44_1(r, p) }
func chain44_1(r Router, p string) { r.Get(p+"/deep1", h); chain44_2(r, p) }
func chain44_2(r Router, p string) { r.Get(p+"/deep2", h); chain44_3(r, p) }
func chain44_3(r Router, p string) { deep_0(r, p) }

type m45 struct{}

func (m45) Register(r Router) { reg45(r, "/m45") }
func reg45(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain45_0(r, p)
}
func chain45_0(r Router, p string) { r.Get(p+"/deep0", h); chain45_1(r, p) }
func chain45_1(r Router, p string) { r.Get(p+"/deep1", h); chain45_2(r, p) }
func chain45_2(r Router, p string) { r.Get(p+"/deep2", h); chain45_3(r, p) }
func chain45_3(r Router, p string) { deep_0(r, p) }

type m46 struct{}

func (m46) Register(r Router) { reg46(r, "/m46") }
func reg46(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain46_0(r, p)
}
func chain46_0(r Router, p string) { r.Get(p+"/deep0", h); chain46_1(r, p) }
func chain46_1(r Router, p string) { r.Get(p+"/deep1", h); chain46_2(r, p) }
func chain46_2(r Router, p string) { r.Get(p+"/deep2", h); chain46_3(r, p) }
func chain46_3(r Router, p string) { deep_0(r, p) }

type m47 struct{}

func (m47) Register(r Router) { reg47(r, "/m47") }
func reg47(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain47_0(r, p)
}
func chain47_0(r Router, p string) { r.Get(p+"/deep0", h); chain47_1(r, p) }
func chain47_1(r Router, p string) { r.Get(p+"/deep1", h); chain47_2(r, p) }
func chain47_2(r Router, p string) { r.Get(p+"/deep2", h); chain47_3(r, p) }
func chain47_3(r Router, p string) { deep_0(r, p) }

type m48 struct{}

func (m48) Register(r Router) { reg48(r, "/m48") }
func reg48(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain48_0(r, p)
}
func chain48_0(r Router, p string) { r.Get(p+"/deep0", h); chain48_1(r, p) }
func chain48_1(r Router, p string) { r.Get(p+"/deep1", h); chain48_2(r, p) }
func chain48_2(r Router, p string) { r.Get(p+"/deep2", h); chain48_3(r, p) }
func chain48_3(r Router, p string) { deep_0(r, p) }

type m49 struct{}

func (m49) Register(r Router) { reg49(r, "/m49") }
func reg49(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain49_0(r, p)
}
func chain49_0(r Router, p string) { r.Get(p+"/deep0", h); chain49_1(r, p) }
func chain49_1(r Router, p string) { r.Get(p+"/deep1", h); chain49_2(r, p) }
func chain49_2(r Router, p string) { r.Get(p+"/deep2", h); chain49_3(r, p) }
func chain49_3(r Router, p string) { deep_0(r, p) }

type m50 struct{}

func (m50) Register(r Router) { reg50(r, "/m50") }
func reg50(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain50_0(r, p)
}
func chain50_0(r Router, p string) { r.Get(p+"/deep0", h); chain50_1(r, p) }
func chain50_1(r Router, p string) { r.Get(p+"/deep1", h); chain50_2(r, p) }
func chain50_2(r Router, p string) { r.Get(p+"/deep2", h); chain50_3(r, p) }
func chain50_3(r Router, p string) { deep_0(r, p) }

type m51 struct{}

func (m51) Register(r Router) { reg51(r, "/m51") }
func reg51(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain51_0(r, p)
}
func chain51_0(r Router, p string) { r.Get(p+"/deep0", h); chain51_1(r, p) }
func chain51_1(r Router, p string) { r.Get(p+"/deep1", h); chain51_2(r, p) }
func chain51_2(r Router, p string) { r.Get(p+"/deep2", h); chain51_3(r, p) }
func chain51_3(r Router, p string) { deep_0(r, p) }

type m52 struct{}

func (m52) Register(r Router) { reg52(r, "/m52") }
func reg52(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain52_0(r, p)
}
func chain52_0(r Router, p string) { r.Get(p+"/deep0", h); chain52_1(r, p) }
func chain52_1(r Router, p string) { r.Get(p+"/deep1", h); chain52_2(r, p) }
func chain52_2(r Router, p string) { r.Get(p+"/deep2", h); chain52_3(r, p) }
func chain52_3(r Router, p string) { deep_0(r, p) }

type m53 struct{}

func (m53) Register(r Router) { reg53(r, "/m53") }
func reg53(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain53_0(r, p)
}
func chain53_0(r Router, p string) { r.Get(p+"/deep0", h); chain53_1(r, p) }
func chain53_1(r Router, p string) { r.Get(p+"/deep1", h); chain53_2(r, p) }
func chain53_2(r Router, p string) { r.Get(p+"/deep2", h); chain53_3(r, p) }
func chain53_3(r Router, p string) { deep_0(r, p) }

type m54 struct{}

func (m54) Register(r Router) { reg54(r, "/m54") }
func reg54(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain54_0(r, p)
}
func chain54_0(r Router, p string) { r.Get(p+"/deep0", h); chain54_1(r, p) }
func chain54_1(r Router, p string) { r.Get(p+"/deep1", h); chain54_2(r, p) }
func chain54_2(r Router, p string) { r.Get(p+"/deep2", h); chain54_3(r, p) }
func chain54_3(r Router, p string) { deep_0(r, p) }

type m55 struct{}

func (m55) Register(r Router) { reg55(r, "/m55") }
func reg55(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain55_0(r, p)
}
func chain55_0(r Router, p string) { r.Get(p+"/deep0", h); chain55_1(r, p) }
func chain55_1(r Router, p string) { r.Get(p+"/deep1", h); chain55_2(r, p) }
func chain55_2(r Router, p string) { r.Get(p+"/deep2", h); chain55_3(r, p) }
func chain55_3(r Router, p string) { deep_0(r, p) }

type m56 struct{}

func (m56) Register(r Router) { reg56(r, "/m56") }
func reg56(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain56_0(r, p)
}
func chain56_0(r Router, p string) { r.Get(p+"/deep0", h); chain56_1(r, p) }
func chain56_1(r Router, p string) { r.Get(p+"/deep1", h); chain56_2(r, p) }
func chain56_2(r Router, p string) { r.Get(p+"/deep2", h); chain56_3(r, p) }
func chain56_3(r Router, p string) { deep_0(r, p) }

type m57 struct{}

func (m57) Register(r Router) { reg57(r, "/m57") }
func reg57(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain57_0(r, p)
}
func chain57_0(r Router, p string) { r.Get(p+"/deep0", h); chain57_1(r, p) }
func chain57_1(r Router, p string) { r.Get(p+"/deep1", h); chain57_2(r, p) }
func chain57_2(r Router, p string) { r.Get(p+"/deep2", h); chain57_3(r, p) }
func chain57_3(r Router, p string) { deep_0(r, p) }

type m58 struct{}

func (m58) Register(r Router) { reg58(r, "/m58") }
func reg58(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain58_0(r, p)
}
func chain58_0(r Router, p string) { r.Get(p+"/deep0", h); chain58_1(r, p) }
func chain58_1(r Router, p string) { r.Get(p+"/deep1", h); chain58_2(r, p) }
func chain58_2(r Router, p string) { r.Get(p+"/deep2", h); chain58_3(r, p) }
func chain58_3(r Router, p string) { deep_0(r, p) }

type m59 struct{}

func (m59) Register(r Router) { reg59(r, "/m59") }
func reg59(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain59_0(r, p)
}
func chain59_0(r Router, p string) { r.Get(p+"/deep0", h); chain59_1(r, p) }
func chain59_1(r Router, p string) { r.Get(p+"/deep1", h); chain59_2(r, p) }
func chain59_2(r Router, p string) { r.Get(p+"/deep2", h); chain59_3(r, p) }
func chain59_3(r Router, p string) { deep_0(r, p) }

type m60 struct{}

func (m60) Register(r Router) { reg60(r, "/m60") }
func reg60(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain60_0(r, p)
}
func chain60_0(r Router, p string) { r.Get(p+"/deep0", h); chain60_1(r, p) }
func chain60_1(r Router, p string) { r.Get(p+"/deep1", h); chain60_2(r, p) }
func chain60_2(r Router, p string) { r.Get(p+"/deep2", h); chain60_3(r, p) }
func chain60_3(r Router, p string) { deep_0(r, p) }

type m61 struct{}

func (m61) Register(r Router) { reg61(r, "/m61") }
func reg61(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain61_0(r, p)
}
func chain61_0(r Router, p string) { r.Get(p+"/deep0", h); chain61_1(r, p) }
func chain61_1(r Router, p string) { r.Get(p+"/deep1", h); chain61_2(r, p) }
func chain61_2(r Router, p string) { r.Get(p+"/deep2", h); chain61_3(r, p) }
func chain61_3(r Router, p string) { deep_0(r, p) }

type m62 struct{}

func (m62) Register(r Router) { reg62(r, "/m62") }
func reg62(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain62_0(r, p)
}
func chain62_0(r Router, p string) { r.Get(p+"/deep0", h); chain62_1(r, p) }
func chain62_1(r Router, p string) { r.Get(p+"/deep1", h); chain62_2(r, p) }
func chain62_2(r Router, p string) { r.Get(p+"/deep2", h); chain62_3(r, p) }
func chain62_3(r Router, p string) { deep_0(r, p) }

type m63 struct{}

func (m63) Register(r Router) { reg63(r, "/m63") }
func reg63(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain63_0(r, p)
}
func chain63_0(r Router, p string) { r.Get(p+"/deep0", h); chain63_1(r, p) }
func chain63_1(r Router, p string) { r.Get(p+"/deep1", h); chain63_2(r, p) }
func chain63_2(r Router, p string) { r.Get(p+"/deep2", h); chain63_3(r, p) }
func chain63_3(r Router, p string) { deep_0(r, p) }

type m64 struct{}

func (m64) Register(r Router) { reg64(r, "/m64") }
func reg64(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain64_0(r, p)
}
func chain64_0(r Router, p string) { r.Get(p+"/deep0", h); chain64_1(r, p) }
func chain64_1(r Router, p string) { r.Get(p+"/deep1", h); chain64_2(r, p) }
func chain64_2(r Router, p string) { r.Get(p+"/deep2", h); chain64_3(r, p) }
func chain64_3(r Router, p string) { deep_0(r, p) }

type m65 struct{}

func (m65) Register(r Router) { reg65(r, "/m65") }
func reg65(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain65_0(r, p)
}
func chain65_0(r Router, p string) { r.Get(p+"/deep0", h); chain65_1(r, p) }
func chain65_1(r Router, p string) { r.Get(p+"/deep1", h); chain65_2(r, p) }
func chain65_2(r Router, p string) { r.Get(p+"/deep2", h); chain65_3(r, p) }
func chain65_3(r Router, p string) { deep_0(r, p) }

type m66 struct{}

func (m66) Register(r Router) { reg66(r, "/m66") }
func reg66(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain66_0(r, p)
}
func chain66_0(r Router, p string) { r.Get(p+"/deep0", h); chain66_1(r, p) }
func chain66_1(r Router, p string) { r.Get(p+"/deep1", h); chain66_2(r, p) }
func chain66_2(r Router, p string) { r.Get(p+"/deep2", h); chain66_3(r, p) }
func chain66_3(r Router, p string) { deep_0(r, p) }

type m67 struct{}

func (m67) Register(r Router) { reg67(r, "/m67") }
func reg67(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain67_0(r, p)
}
func chain67_0(r Router, p string) { r.Get(p+"/deep0", h); chain67_1(r, p) }
func chain67_1(r Router, p string) { r.Get(p+"/deep1", h); chain67_2(r, p) }
func chain67_2(r Router, p string) { r.Get(p+"/deep2", h); chain67_3(r, p) }
func chain67_3(r Router, p string) { deep_0(r, p) }

type m68 struct{}

func (m68) Register(r Router) { reg68(r, "/m68") }
func reg68(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain68_0(r, p)
}
func chain68_0(r Router, p string) { r.Get(p+"/deep0", h); chain68_1(r, p) }
func chain68_1(r Router, p string) { r.Get(p+"/deep1", h); chain68_2(r, p) }
func chain68_2(r Router, p string) { r.Get(p+"/deep2", h); chain68_3(r, p) }
func chain68_3(r Router, p string) { deep_0(r, p) }

type m69 struct{}

func (m69) Register(r Router) { reg69(r, "/m69") }
func reg69(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain69_0(r, p)
}
func chain69_0(r Router, p string) { r.Get(p+"/deep0", h); chain69_1(r, p) }
func chain69_1(r Router, p string) { r.Get(p+"/deep1", h); chain69_2(r, p) }
func chain69_2(r Router, p string) { r.Get(p+"/deep2", h); chain69_3(r, p) }
func chain69_3(r Router, p string) { deep_0(r, p) }

type m70 struct{}

func (m70) Register(r Router) { reg70(r, "/m70") }
func reg70(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain70_0(r, p)
}
func chain70_0(r Router, p string) { r.Get(p+"/deep0", h); chain70_1(r, p) }
func chain70_1(r Router, p string) { r.Get(p+"/deep1", h); chain70_2(r, p) }
func chain70_2(r Router, p string) { r.Get(p+"/deep2", h); chain70_3(r, p) }
func chain70_3(r Router, p string) { deep_0(r, p) }

type m71 struct{}

func (m71) Register(r Router) { reg71(r, "/m71") }
func reg71(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain71_0(r, p)
}
func chain71_0(r Router, p string) { r.Get(p+"/deep0", h); chain71_1(r, p) }
func chain71_1(r Router, p string) { r.Get(p+"/deep1", h); chain71_2(r, p) }
func chain71_2(r Router, p string) { r.Get(p+"/deep2", h); chain71_3(r, p) }
func chain71_3(r Router, p string) { deep_0(r, p) }

type m72 struct{}

func (m72) Register(r Router) { reg72(r, "/m72") }
func reg72(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain72_0(r, p)
}
func chain72_0(r Router, p string) { r.Get(p+"/deep0", h); chain72_1(r, p) }
func chain72_1(r Router, p string) { r.Get(p+"/deep1", h); chain72_2(r, p) }
func chain72_2(r Router, p string) { r.Get(p+"/deep2", h); chain72_3(r, p) }
func chain72_3(r Router, p string) { deep_0(r, p) }

type m73 struct{}

func (m73) Register(r Router) { reg73(r, "/m73") }
func reg73(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain73_0(r, p)
}
func chain73_0(r Router, p string) { r.Get(p+"/deep0", h); chain73_1(r, p) }
func chain73_1(r Router, p string) { r.Get(p+"/deep1", h); chain73_2(r, p) }
func chain73_2(r Router, p string) { r.Get(p+"/deep2", h); chain73_3(r, p) }
func chain73_3(r Router, p string) { deep_0(r, p) }

type m74 struct{}

func (m74) Register(r Router) { reg74(r, "/m74") }
func reg74(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain74_0(r, p)
}
func chain74_0(r Router, p string) { r.Get(p+"/deep0", h); chain74_1(r, p) }
func chain74_1(r Router, p string) { r.Get(p+"/deep1", h); chain74_2(r, p) }
func chain74_2(r Router, p string) { r.Get(p+"/deep2", h); chain74_3(r, p) }
func chain74_3(r Router, p string) { deep_0(r, p) }

type m75 struct{}

func (m75) Register(r Router) { reg75(r, "/m75") }
func reg75(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain75_0(r, p)
}
func chain75_0(r Router, p string) { r.Get(p+"/deep0", h); chain75_1(r, p) }
func chain75_1(r Router, p string) { r.Get(p+"/deep1", h); chain75_2(r, p) }
func chain75_2(r Router, p string) { r.Get(p+"/deep2", h); chain75_3(r, p) }
func chain75_3(r Router, p string) { deep_0(r, p) }

type m76 struct{}

func (m76) Register(r Router) { reg76(r, "/m76") }
func reg76(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain76_0(r, p)
}
func chain76_0(r Router, p string) { r.Get(p+"/deep0", h); chain76_1(r, p) }
func chain76_1(r Router, p string) { r.Get(p+"/deep1", h); chain76_2(r, p) }
func chain76_2(r Router, p string) { r.Get(p+"/deep2", h); chain76_3(r, p) }
func chain76_3(r Router, p string) { deep_0(r, p) }

type m77 struct{}

func (m77) Register(r Router) { reg77(r, "/m77") }
func reg77(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain77_0(r, p)
}
func chain77_0(r Router, p string) { r.Get(p+"/deep0", h); chain77_1(r, p) }
func chain77_1(r Router, p string) { r.Get(p+"/deep1", h); chain77_2(r, p) }
func chain77_2(r Router, p string) { r.Get(p+"/deep2", h); chain77_3(r, p) }
func chain77_3(r Router, p string) { deep_0(r, p) }

type m78 struct{}

func (m78) Register(r Router) { reg78(r, "/m78") }
func reg78(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain78_0(r, p)
}
func chain78_0(r Router, p string) { r.Get(p+"/deep0", h); chain78_1(r, p) }
func chain78_1(r Router, p string) { r.Get(p+"/deep1", h); chain78_2(r, p) }
func chain78_2(r Router, p string) { r.Get(p+"/deep2", h); chain78_3(r, p) }
func chain78_3(r Router, p string) { deep_0(r, p) }

type m79 struct{}

func (m79) Register(r Router) { reg79(r, "/m79") }
func reg79(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain79_0(r, p)
}
func chain79_0(r Router, p string) { r.Get(p+"/deep0", h); chain79_1(r, p) }
func chain79_1(r Router, p string) { r.Get(p+"/deep1", h); chain79_2(r, p) }
func chain79_2(r Router, p string) { r.Get(p+"/deep2", h); chain79_3(r, p) }
func chain79_3(r Router, p string) { deep_0(r, p) }

type m80 struct{}

func (m80) Register(r Router) { reg80(r, "/m80") }
func reg80(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain80_0(r, p)
}
func chain80_0(r Router, p string) { r.Get(p+"/deep0", h); chain80_1(r, p) }
func chain80_1(r Router, p string) { r.Get(p+"/deep1", h); chain80_2(r, p) }
func chain80_2(r Router, p string) { r.Get(p+"/deep2", h); chain80_3(r, p) }
func chain80_3(r Router, p string) { deep_0(r, p) }

type m81 struct{}

func (m81) Register(r Router) { reg81(r, "/m81") }
func reg81(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain81_0(r, p)
}
func chain81_0(r Router, p string) { r.Get(p+"/deep0", h); chain81_1(r, p) }
func chain81_1(r Router, p string) { r.Get(p+"/deep1", h); chain81_2(r, p) }
func chain81_2(r Router, p string) { r.Get(p+"/deep2", h); chain81_3(r, p) }
func chain81_3(r Router, p string) { deep_0(r, p) }

type m82 struct{}

func (m82) Register(r Router) { reg82(r, "/m82") }
func reg82(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain82_0(r, p)
}
func chain82_0(r Router, p string) { r.Get(p+"/deep0", h); chain82_1(r, p) }
func chain82_1(r Router, p string) { r.Get(p+"/deep1", h); chain82_2(r, p) }
func chain82_2(r Router, p string) { r.Get(p+"/deep2", h); chain82_3(r, p) }
func chain82_3(r Router, p string) { deep_0(r, p) }

type m83 struct{}

func (m83) Register(r Router) { reg83(r, "/m83") }
func reg83(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain83_0(r, p)
}
func chain83_0(r Router, p string) { r.Get(p+"/deep0", h); chain83_1(r, p) }
func chain83_1(r Router, p string) { r.Get(p+"/deep1", h); chain83_2(r, p) }
func chain83_2(r Router, p string) { r.Get(p+"/deep2", h); chain83_3(r, p) }
func chain83_3(r Router, p string) { deep_0(r, p) }

type m84 struct{}

func (m84) Register(r Router) { reg84(r, "/m84") }
func reg84(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain84_0(r, p)
}
func chain84_0(r Router, p string) { r.Get(p+"/deep0", h); chain84_1(r, p) }
func chain84_1(r Router, p string) { r.Get(p+"/deep1", h); chain84_2(r, p) }
func chain84_2(r Router, p string) { r.Get(p+"/deep2", h); chain84_3(r, p) }
func chain84_3(r Router, p string) { deep_0(r, p) }

type m85 struct{}

func (m85) Register(r Router) { reg85(r, "/m85") }
func reg85(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain85_0(r, p)
}
func chain85_0(r Router, p string) { r.Get(p+"/deep0", h); chain85_1(r, p) }
func chain85_1(r Router, p string) { r.Get(p+"/deep1", h); chain85_2(r, p) }
func chain85_2(r Router, p string) { r.Get(p+"/deep2", h); chain85_3(r, p) }
func chain85_3(r Router, p string) { deep_0(r, p) }

type m86 struct{}

func (m86) Register(r Router) { reg86(r, "/m86") }
func reg86(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain86_0(r, p)
}
func chain86_0(r Router, p string) { r.Get(p+"/deep0", h); chain86_1(r, p) }
func chain86_1(r Router, p string) { r.Get(p+"/deep1", h); chain86_2(r, p) }
func chain86_2(r Router, p string) { r.Get(p+"/deep2", h); chain86_3(r, p) }
func chain86_3(r Router, p string) { deep_0(r, p) }

type m87 struct{}

func (m87) Register(r Router) { reg87(r, "/m87") }
func reg87(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain87_0(r, p)
}
func chain87_0(r Router, p string) { r.Get(p+"/deep0", h); chain87_1(r, p) }
func chain87_1(r Router, p string) { r.Get(p+"/deep1", h); chain87_2(r, p) }
func chain87_2(r Router, p string) { r.Get(p+"/deep2", h); chain87_3(r, p) }
func chain87_3(r Router, p string) { deep_0(r, p) }

type m88 struct{}

func (m88) Register(r Router) { reg88(r, "/m88") }
func reg88(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain88_0(r, p)
}
func chain88_0(r Router, p string) { r.Get(p+"/deep0", h); chain88_1(r, p) }
func chain88_1(r Router, p string) { r.Get(p+"/deep1", h); chain88_2(r, p) }
func chain88_2(r Router, p string) { r.Get(p+"/deep2", h); chain88_3(r, p) }
func chain88_3(r Router, p string) { deep_0(r, p) }

type m89 struct{}

func (m89) Register(r Router) { reg89(r, "/m89") }
func reg89(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain89_0(r, p)
}
func chain89_0(r Router, p string) { r.Get(p+"/deep0", h); chain89_1(r, p) }
func chain89_1(r Router, p string) { r.Get(p+"/deep1", h); chain89_2(r, p) }
func chain89_2(r Router, p string) { r.Get(p+"/deep2", h); chain89_3(r, p) }
func chain89_3(r Router, p string) { deep_0(r, p) }

type m90 struct{}

func (m90) Register(r Router) { reg90(r, "/m90") }
func reg90(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain90_0(r, p)
}
func chain90_0(r Router, p string) { r.Get(p+"/deep0", h); chain90_1(r, p) }
func chain90_1(r Router, p string) { r.Get(p+"/deep1", h); chain90_2(r, p) }
func chain90_2(r Router, p string) { r.Get(p+"/deep2", h); chain90_3(r, p) }
func chain90_3(r Router, p string) { deep_0(r, p) }

type m91 struct{}

func (m91) Register(r Router) { reg91(r, "/m91") }
func reg91(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain91_0(r, p)
}
func chain91_0(r Router, p string) { r.Get(p+"/deep0", h); chain91_1(r, p) }
func chain91_1(r Router, p string) { r.Get(p+"/deep1", h); chain91_2(r, p) }
func chain91_2(r Router, p string) { r.Get(p+"/deep2", h); chain91_3(r, p) }
func chain91_3(r Router, p string) { deep_0(r, p) }

type m92 struct{}

func (m92) Register(r Router) { reg92(r, "/m92") }
func reg92(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain92_0(r, p)
}
func chain92_0(r Router, p string) { r.Get(p+"/deep0", h); chain92_1(r, p) }
func chain92_1(r Router, p string) { r.Get(p+"/deep1", h); chain92_2(r, p) }
func chain92_2(r Router, p string) { r.Get(p+"/deep2", h); chain92_3(r, p) }
func chain92_3(r Router, p string) { deep_0(r, p) }

type m93 struct{}

func (m93) Register(r Router) { reg93(r, "/m93") }
func reg93(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain93_0(r, p)
}
func chain93_0(r Router, p string) { r.Get(p+"/deep0", h); chain93_1(r, p) }
func chain93_1(r Router, p string) { r.Get(p+"/deep1", h); chain93_2(r, p) }
func chain93_2(r Router, p string) { r.Get(p+"/deep2", h); chain93_3(r, p) }
func chain93_3(r Router, p string) { deep_0(r, p) }

type m94 struct{}

func (m94) Register(r Router) { reg94(r, "/m94") }
func reg94(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain94_0(r, p)
}
func chain94_0(r Router, p string) { r.Get(p+"/deep0", h); chain94_1(r, p) }
func chain94_1(r Router, p string) { r.Get(p+"/deep1", h); chain94_2(r, p) }
func chain94_2(r Router, p string) { r.Get(p+"/deep2", h); chain94_3(r, p) }
func chain94_3(r Router, p string) { deep_0(r, p) }

type m95 struct{}

func (m95) Register(r Router) { reg95(r, "/m95") }
func reg95(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain95_0(r, p)
}
func chain95_0(r Router, p string) { r.Get(p+"/deep0", h); chain95_1(r, p) }
func chain95_1(r Router, p string) { r.Get(p+"/deep1", h); chain95_2(r, p) }
func chain95_2(r Router, p string) { r.Get(p+"/deep2", h); chain95_3(r, p) }
func chain95_3(r Router, p string) { deep_0(r, p) }

type m96 struct{}

func (m96) Register(r Router) { reg96(r, "/m96") }
func reg96(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain96_0(r, p)
}
func chain96_0(r Router, p string) { r.Get(p+"/deep0", h); chain96_1(r, p) }
func chain96_1(r Router, p string) { r.Get(p+"/deep1", h); chain96_2(r, p) }
func chain96_2(r Router, p string) { r.Get(p+"/deep2", h); chain96_3(r, p) }
func chain96_3(r Router, p string) { deep_0(r, p) }

type m97 struct{}

func (m97) Register(r Router) { reg97(r, "/m97") }
func reg97(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain97_0(r, p)
}
func chain97_0(r Router, p string) { r.Get(p+"/deep0", h); chain97_1(r, p) }
func chain97_1(r Router, p string) { r.Get(p+"/deep1", h); chain97_2(r, p) }
func chain97_2(r Router, p string) { r.Get(p+"/deep2", h); chain97_3(r, p) }
func chain97_3(r Router, p string) { deep_0(r, p) }

type m98 struct{}

func (m98) Register(r Router) { reg98(r, "/m98") }
func reg98(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain98_0(r, p)
}
func chain98_0(r Router, p string) { r.Get(p+"/deep0", h); chain98_1(r, p) }
func chain98_1(r Router, p string) { r.Get(p+"/deep1", h); chain98_2(r, p) }
func chain98_2(r Router, p string) { r.Get(p+"/deep2", h); chain98_3(r, p) }
func chain98_3(r Router, p string) { deep_0(r, p) }

type m99 struct{}

func (m99) Register(r Router) { reg99(r, "/m99") }
func reg99(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain99_0(r, p)
}
func chain99_0(r Router, p string) { r.Get(p+"/deep0", h); chain99_1(r, p) }
func chain99_1(r Router, p string) { r.Get(p+"/deep1", h); chain99_2(r, p) }
func chain99_2(r Router, p string) { r.Get(p+"/deep2", h); chain99_3(r, p) }
func chain99_3(r Router, p string) { deep_0(r, p) }

type m100 struct{}

func (m100) Register(r Router) { reg100(r, "/m100") }
func reg100(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain100_0(r, p)
}
func chain100_0(r Router, p string) { r.Get(p+"/deep0", h); chain100_1(r, p) }
func chain100_1(r Router, p string) { r.Get(p+"/deep1", h); chain100_2(r, p) }
func chain100_2(r Router, p string) { r.Get(p+"/deep2", h); chain100_3(r, p) }
func chain100_3(r Router, p string) { deep_0(r, p) }

type m101 struct{}

func (m101) Register(r Router) { reg101(r, "/m101") }
func reg101(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain101_0(r, p)
}
func chain101_0(r Router, p string) { r.Get(p+"/deep0", h); chain101_1(r, p) }
func chain101_1(r Router, p string) { r.Get(p+"/deep1", h); chain101_2(r, p) }
func chain101_2(r Router, p string) { r.Get(p+"/deep2", h); chain101_3(r, p) }
func chain101_3(r Router, p string) { deep_0(r, p) }

type m102 struct{}

func (m102) Register(r Router) { reg102(r, "/m102") }
func reg102(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain102_0(r, p)
}
func chain102_0(r Router, p string) { r.Get(p+"/deep0", h); chain102_1(r, p) }
func chain102_1(r Router, p string) { r.Get(p+"/deep1", h); chain102_2(r, p) }
func chain102_2(r Router, p string) { r.Get(p+"/deep2", h); chain102_3(r, p) }
func chain102_3(r Router, p string) { deep_0(r, p) }

type m103 struct{}

func (m103) Register(r Router) { reg103(r, "/m103") }
func reg103(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain103_0(r, p)
}
func chain103_0(r Router, p string) { r.Get(p+"/deep0", h); chain103_1(r, p) }
func chain103_1(r Router, p string) { r.Get(p+"/deep1", h); chain103_2(r, p) }
func chain103_2(r Router, p string) { r.Get(p+"/deep2", h); chain103_3(r, p) }
func chain103_3(r Router, p string) { deep_0(r, p) }

type m104 struct{}

func (m104) Register(r Router) { reg104(r, "/m104") }
func reg104(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain104_0(r, p)
}
func chain104_0(r Router, p string) { r.Get(p+"/deep0", h); chain104_1(r, p) }
func chain104_1(r Router, p string) { r.Get(p+"/deep1", h); chain104_2(r, p) }
func chain104_2(r Router, p string) { r.Get(p+"/deep2", h); chain104_3(r, p) }
func chain104_3(r Router, p string) { deep_0(r, p) }

type m105 struct{}

func (m105) Register(r Router) { reg105(r, "/m105") }
func reg105(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain105_0(r, p)
}
func chain105_0(r Router, p string) { r.Get(p+"/deep0", h); chain105_1(r, p) }
func chain105_1(r Router, p string) { r.Get(p+"/deep1", h); chain105_2(r, p) }
func chain105_2(r Router, p string) { r.Get(p+"/deep2", h); chain105_3(r, p) }
func chain105_3(r Router, p string) { deep_0(r, p) }

type m106 struct{}

func (m106) Register(r Router) { reg106(r, "/m106") }
func reg106(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain106_0(r, p)
}
func chain106_0(r Router, p string) { r.Get(p+"/deep0", h); chain106_1(r, p) }
func chain106_1(r Router, p string) { r.Get(p+"/deep1", h); chain106_2(r, p) }
func chain106_2(r Router, p string) { r.Get(p+"/deep2", h); chain106_3(r, p) }
func chain106_3(r Router, p string) { deep_0(r, p) }

type m107 struct{}

func (m107) Register(r Router) { reg107(r, "/m107") }
func reg107(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain107_0(r, p)
}
func chain107_0(r Router, p string) { r.Get(p+"/deep0", h); chain107_1(r, p) }
func chain107_1(r Router, p string) { r.Get(p+"/deep1", h); chain107_2(r, p) }
func chain107_2(r Router, p string) { r.Get(p+"/deep2", h); chain107_3(r, p) }
func chain107_3(r Router, p string) { deep_0(r, p) }

type m108 struct{}

func (m108) Register(r Router) { reg108(r, "/m108") }
func reg108(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain108_0(r, p)
}
func chain108_0(r Router, p string) { r.Get(p+"/deep0", h); chain108_1(r, p) }
func chain108_1(r Router, p string) { r.Get(p+"/deep1", h); chain108_2(r, p) }
func chain108_2(r Router, p string) { r.Get(p+"/deep2", h); chain108_3(r, p) }
func chain108_3(r Router, p string) { deep_0(r, p) }

type m109 struct{}

func (m109) Register(r Router) { reg109(r, "/m109") }
func reg109(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain109_0(r, p)
}
func chain109_0(r Router, p string) { r.Get(p+"/deep0", h); chain109_1(r, p) }
func chain109_1(r Router, p string) { r.Get(p+"/deep1", h); chain109_2(r, p) }
func chain109_2(r Router, p string) { r.Get(p+"/deep2", h); chain109_3(r, p) }
func chain109_3(r Router, p string) { deep_0(r, p) }

type m110 struct{}

func (m110) Register(r Router) { reg110(r, "/m110") }
func reg110(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain110_0(r, p)
}
func chain110_0(r Router, p string) { r.Get(p+"/deep0", h); chain110_1(r, p) }
func chain110_1(r Router, p string) { r.Get(p+"/deep1", h); chain110_2(r, p) }
func chain110_2(r Router, p string) { r.Get(p+"/deep2", h); chain110_3(r, p) }
func chain110_3(r Router, p string) { deep_0(r, p) }

type m111 struct{}

func (m111) Register(r Router) { reg111(r, "/m111") }
func reg111(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain111_0(r, p)
}
func chain111_0(r Router, p string) { r.Get(p+"/deep0", h); chain111_1(r, p) }
func chain111_1(r Router, p string) { r.Get(p+"/deep1", h); chain111_2(r, p) }
func chain111_2(r Router, p string) { r.Get(p+"/deep2", h); chain111_3(r, p) }
func chain111_3(r Router, p string) { deep_0(r, p) }

type m112 struct{}

func (m112) Register(r Router) { reg112(r, "/m112") }
func reg112(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain112_0(r, p)
}
func chain112_0(r Router, p string) { r.Get(p+"/deep0", h); chain112_1(r, p) }
func chain112_1(r Router, p string) { r.Get(p+"/deep1", h); chain112_2(r, p) }
func chain112_2(r Router, p string) { r.Get(p+"/deep2", h); chain112_3(r, p) }
func chain112_3(r Router, p string) { deep_0(r, p) }

type m113 struct{}

func (m113) Register(r Router) { reg113(r, "/m113") }
func reg113(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain113_0(r, p)
}
func chain113_0(r Router, p string) { r.Get(p+"/deep0", h); chain113_1(r, p) }
func chain113_1(r Router, p string) { r.Get(p+"/deep1", h); chain113_2(r, p) }
func chain113_2(r Router, p string) { r.Get(p+"/deep2", h); chain113_3(r, p) }
func chain113_3(r Router, p string) { deep_0(r, p) }

type m114 struct{}

func (m114) Register(r Router) { reg114(r, "/m114") }
func reg114(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain114_0(r, p)
}
func chain114_0(r Router, p string) { r.Get(p+"/deep0", h); chain114_1(r, p) }
func chain114_1(r Router, p string) { r.Get(p+"/deep1", h); chain114_2(r, p) }
func chain114_2(r Router, p string) { r.Get(p+"/deep2", h); chain114_3(r, p) }
func chain114_3(r Router, p string) { deep_0(r, p) }

type m115 struct{}

func (m115) Register(r Router) { reg115(r, "/m115") }
func reg115(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain115_0(r, p)
}
func chain115_0(r Router, p string) { r.Get(p+"/deep0", h); chain115_1(r, p) }
func chain115_1(r Router, p string) { r.Get(p+"/deep1", h); chain115_2(r, p) }
func chain115_2(r Router, p string) { r.Get(p+"/deep2", h); chain115_3(r, p) }
func chain115_3(r Router, p string) { deep_0(r, p) }

type m116 struct{}

func (m116) Register(r Router) { reg116(r, "/m116") }
func reg116(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain116_0(r, p)
}
func chain116_0(r Router, p string) { r.Get(p+"/deep0", h); chain116_1(r, p) }
func chain116_1(r Router, p string) { r.Get(p+"/deep1", h); chain116_2(r, p) }
func chain116_2(r Router, p string) { r.Get(p+"/deep2", h); chain116_3(r, p) }
func chain116_3(r Router, p string) { deep_0(r, p) }

type m117 struct{}

func (m117) Register(r Router) { reg117(r, "/m117") }
func reg117(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain117_0(r, p)
}
func chain117_0(r Router, p string) { r.Get(p+"/deep0", h); chain117_1(r, p) }
func chain117_1(r Router, p string) { r.Get(p+"/deep1", h); chain117_2(r, p) }
func chain117_2(r Router, p string) { r.Get(p+"/deep2", h); chain117_3(r, p) }
func chain117_3(r Router, p string) { deep_0(r, p) }

type m118 struct{}

func (m118) Register(r Router) { reg118(r, "/m118") }
func reg118(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain118_0(r, p)
}
func chain118_0(r Router, p string) { r.Get(p+"/deep0", h); chain118_1(r, p) }
func chain118_1(r Router, p string) { r.Get(p+"/deep1", h); chain118_2(r, p) }
func chain118_2(r Router, p string) { r.Get(p+"/deep2", h); chain118_3(r, p) }
func chain118_3(r Router, p string) { deep_0(r, p) }

type m119 struct{}

func (m119) Register(r Router) { reg119(r, "/m119") }
func reg119(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain119_0(r, p)
}
func chain119_0(r Router, p string) { r.Get(p+"/deep0", h); chain119_1(r, p) }
func chain119_1(r Router, p string) { r.Get(p+"/deep1", h); chain119_2(r, p) }
func chain119_2(r Router, p string) { r.Get(p+"/deep2", h); chain119_3(r, p) }
func chain119_3(r Router, p string) { deep_0(r, p) }

type m120 struct{}

func (m120) Register(r Router) { reg120(r, "/m120") }
func reg120(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain120_0(r, p)
}
func chain120_0(r Router, p string) { r.Get(p+"/deep0", h); chain120_1(r, p) }
func chain120_1(r Router, p string) { r.Get(p+"/deep1", h); chain120_2(r, p) }
func chain120_2(r Router, p string) { r.Get(p+"/deep2", h); chain120_3(r, p) }
func chain120_3(r Router, p string) { deep_0(r, p) }

type m121 struct{}

func (m121) Register(r Router) { reg121(r, "/m121") }
func reg121(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain121_0(r, p)
}
func chain121_0(r Router, p string) { r.Get(p+"/deep0", h); chain121_1(r, p) }
func chain121_1(r Router, p string) { r.Get(p+"/deep1", h); chain121_2(r, p) }
func chain121_2(r Router, p string) { r.Get(p+"/deep2", h); chain121_3(r, p) }
func chain121_3(r Router, p string) { deep_0(r, p) }

type m122 struct{}

func (m122) Register(r Router) { reg122(r, "/m122") }
func reg122(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain122_0(r, p)
}
func chain122_0(r Router, p string) { r.Get(p+"/deep0", h); chain122_1(r, p) }
func chain122_1(r Router, p string) { r.Get(p+"/deep1", h); chain122_2(r, p) }
func chain122_2(r Router, p string) { r.Get(p+"/deep2", h); chain122_3(r, p) }
func chain122_3(r Router, p string) { deep_0(r, p) }

type m123 struct{}

func (m123) Register(r Router) { reg123(r, "/m123") }
func reg123(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain123_0(r, p)
}
func chain123_0(r Router, p string) { r.Get(p+"/deep0", h); chain123_1(r, p) }
func chain123_1(r Router, p string) { r.Get(p+"/deep1", h); chain123_2(r, p) }
func chain123_2(r Router, p string) { r.Get(p+"/deep2", h); chain123_3(r, p) }
func chain123_3(r Router, p string) { deep_0(r, p) }

type m124 struct{}

func (m124) Register(r Router) { reg124(r, "/m124") }
func reg124(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain124_0(r, p)
}
func chain124_0(r Router, p string) { r.Get(p+"/deep0", h); chain124_1(r, p) }
func chain124_1(r Router, p string) { r.Get(p+"/deep1", h); chain124_2(r, p) }
func chain124_2(r Router, p string) { r.Get(p+"/deep2", h); chain124_3(r, p) }
func chain124_3(r Router, p string) { deep_0(r, p) }

type m125 struct{}

func (m125) Register(r Router) { reg125(r, "/m125") }
func reg125(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain125_0(r, p)
}
func chain125_0(r Router, p string) { r.Get(p+"/deep0", h); chain125_1(r, p) }
func chain125_1(r Router, p string) { r.Get(p+"/deep1", h); chain125_2(r, p) }
func chain125_2(r Router, p string) { r.Get(p+"/deep2", h); chain125_3(r, p) }
func chain125_3(r Router, p string) { deep_0(r, p) }

type m126 struct{}

func (m126) Register(r Router) { reg126(r, "/m126") }
func reg126(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain126_0(r, p)
}
func chain126_0(r Router, p string) { r.Get(p+"/deep0", h); chain126_1(r, p) }
func chain126_1(r Router, p string) { r.Get(p+"/deep1", h); chain126_2(r, p) }
func chain126_2(r Router, p string) { r.Get(p+"/deep2", h); chain126_3(r, p) }
func chain126_3(r Router, p string) { deep_0(r, p) }

type m127 struct{}

func (m127) Register(r Router) { reg127(r, "/m127") }
func reg127(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain127_0(r, p)
}
func chain127_0(r Router, p string) { r.Get(p+"/deep0", h); chain127_1(r, p) }
func chain127_1(r Router, p string) { r.Get(p+"/deep1", h); chain127_2(r, p) }
func chain127_2(r Router, p string) { r.Get(p+"/deep2", h); chain127_3(r, p) }
func chain127_3(r Router, p string) { deep_0(r, p) }

type m128 struct{}

func (m128) Register(r Router) { reg128(r, "/m128") }
func reg128(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain128_0(r, p)
}
func chain128_0(r Router, p string) { r.Get(p+"/deep0", h); chain128_1(r, p) }
func chain128_1(r Router, p string) { r.Get(p+"/deep1", h); chain128_2(r, p) }
func chain128_2(r Router, p string) { r.Get(p+"/deep2", h); chain128_3(r, p) }
func chain128_3(r Router, p string) { deep_0(r, p) }

type m129 struct{}

func (m129) Register(r Router) { reg129(r, "/m129") }
func reg129(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain129_0(r, p)
}
func chain129_0(r Router, p string) { r.Get(p+"/deep0", h); chain129_1(r, p) }
func chain129_1(r Router, p string) { r.Get(p+"/deep1", h); chain129_2(r, p) }
func chain129_2(r Router, p string) { r.Get(p+"/deep2", h); chain129_3(r, p) }
func chain129_3(r Router, p string) { deep_0(r, p) }

type m130 struct{}

func (m130) Register(r Router) { reg130(r, "/m130") }
func reg130(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain130_0(r, p)
}
func chain130_0(r Router, p string) { r.Get(p+"/deep0", h); chain130_1(r, p) }
func chain130_1(r Router, p string) { r.Get(p+"/deep1", h); chain130_2(r, p) }
func chain130_2(r Router, p string) { r.Get(p+"/deep2", h); chain130_3(r, p) }
func chain130_3(r Router, p string) { deep_0(r, p) }

type m131 struct{}

func (m131) Register(r Router) { reg131(r, "/m131") }
func reg131(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain131_0(r, p)
}
func chain131_0(r Router, p string) { r.Get(p+"/deep0", h); chain131_1(r, p) }
func chain131_1(r Router, p string) { r.Get(p+"/deep1", h); chain131_2(r, p) }
func chain131_2(r Router, p string) { r.Get(p+"/deep2", h); chain131_3(r, p) }
func chain131_3(r Router, p string) { deep_0(r, p) }

type m132 struct{}

func (m132) Register(r Router) { reg132(r, "/m132") }
func reg132(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain132_0(r, p)
}
func chain132_0(r Router, p string) { r.Get(p+"/deep0", h); chain132_1(r, p) }
func chain132_1(r Router, p string) { r.Get(p+"/deep1", h); chain132_2(r, p) }
func chain132_2(r Router, p string) { r.Get(p+"/deep2", h); chain132_3(r, p) }
func chain132_3(r Router, p string) { deep_0(r, p) }

type m133 struct{}

func (m133) Register(r Router) { reg133(r, "/m133") }
func reg133(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain133_0(r, p)
}
func chain133_0(r Router, p string) { r.Get(p+"/deep0", h); chain133_1(r, p) }
func chain133_1(r Router, p string) { r.Get(p+"/deep1", h); chain133_2(r, p) }
func chain133_2(r Router, p string) { r.Get(p+"/deep2", h); chain133_3(r, p) }
func chain133_3(r Router, p string) { deep_0(r, p) }

type m134 struct{}

func (m134) Register(r Router) { reg134(r, "/m134") }
func reg134(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain134_0(r, p)
}
func chain134_0(r Router, p string) { r.Get(p+"/deep0", h); chain134_1(r, p) }
func chain134_1(r Router, p string) { r.Get(p+"/deep1", h); chain134_2(r, p) }
func chain134_2(r Router, p string) { r.Get(p+"/deep2", h); chain134_3(r, p) }
func chain134_3(r Router, p string) { deep_0(r, p) }

type m135 struct{}

func (m135) Register(r Router) { reg135(r, "/m135") }
func reg135(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain135_0(r, p)
}
func chain135_0(r Router, p string) { r.Get(p+"/deep0", h); chain135_1(r, p) }
func chain135_1(r Router, p string) { r.Get(p+"/deep1", h); chain135_2(r, p) }
func chain135_2(r Router, p string) { r.Get(p+"/deep2", h); chain135_3(r, p) }
func chain135_3(r Router, p string) { deep_0(r, p) }

type m136 struct{}

func (m136) Register(r Router) { reg136(r, "/m136") }
func reg136(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain136_0(r, p)
}
func chain136_0(r Router, p string) { r.Get(p+"/deep0", h); chain136_1(r, p) }
func chain136_1(r Router, p string) { r.Get(p+"/deep1", h); chain136_2(r, p) }
func chain136_2(r Router, p string) { r.Get(p+"/deep2", h); chain136_3(r, p) }
func chain136_3(r Router, p string) { deep_0(r, p) }

type m137 struct{}

func (m137) Register(r Router) { reg137(r, "/m137") }
func reg137(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain137_0(r, p)
}
func chain137_0(r Router, p string) { r.Get(p+"/deep0", h); chain137_1(r, p) }
func chain137_1(r Router, p string) { r.Get(p+"/deep1", h); chain137_2(r, p) }
func chain137_2(r Router, p string) { r.Get(p+"/deep2", h); chain137_3(r, p) }
func chain137_3(r Router, p string) { deep_0(r, p) }

type m138 struct{}

func (m138) Register(r Router) { reg138(r, "/m138") }
func reg138(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain138_0(r, p)
}
func chain138_0(r Router, p string) { r.Get(p+"/deep0", h); chain138_1(r, p) }
func chain138_1(r Router, p string) { r.Get(p+"/deep1", h); chain138_2(r, p) }
func chain138_2(r Router, p string) { r.Get(p+"/deep2", h); chain138_3(r, p) }
func chain138_3(r Router, p string) { deep_0(r, p) }

type m139 struct{}

func (m139) Register(r Router) { reg139(r, "/m139") }
func reg139(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain139_0(r, p)
}
func chain139_0(r Router, p string) { r.Get(p+"/deep0", h); chain139_1(r, p) }
func chain139_1(r Router, p string) { r.Get(p+"/deep1", h); chain139_2(r, p) }
func chain139_2(r Router, p string) { r.Get(p+"/deep2", h); chain139_3(r, p) }
func chain139_3(r Router, p string) { deep_0(r, p) }

type m140 struct{}

func (m140) Register(r Router) { reg140(r, "/m140") }
func reg140(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain140_0(r, p)
}
func chain140_0(r Router, p string) { r.Get(p+"/deep0", h); chain140_1(r, p) }
func chain140_1(r Router, p string) { r.Get(p+"/deep1", h); chain140_2(r, p) }
func chain140_2(r Router, p string) { r.Get(p+"/deep2", h); chain140_3(r, p) }
func chain140_3(r Router, p string) { deep_0(r, p) }

type m141 struct{}

func (m141) Register(r Router) { reg141(r, "/m141") }
func reg141(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain141_0(r, p)
}
func chain141_0(r Router, p string) { r.Get(p+"/deep0", h); chain141_1(r, p) }
func chain141_1(r Router, p string) { r.Get(p+"/deep1", h); chain141_2(r, p) }
func chain141_2(r Router, p string) { r.Get(p+"/deep2", h); chain141_3(r, p) }
func chain141_3(r Router, p string) { deep_0(r, p) }

type m142 struct{}

func (m142) Register(r Router) { reg142(r, "/m142") }
func reg142(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain142_0(r, p)
}
func chain142_0(r Router, p string) { r.Get(p+"/deep0", h); chain142_1(r, p) }
func chain142_1(r Router, p string) { r.Get(p+"/deep1", h); chain142_2(r, p) }
func chain142_2(r Router, p string) { r.Get(p+"/deep2", h); chain142_3(r, p) }
func chain142_3(r Router, p string) { deep_0(r, p) }

type m143 struct{}

func (m143) Register(r Router) { reg143(r, "/m143") }
func reg143(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain143_0(r, p)
}
func chain143_0(r Router, p string) { r.Get(p+"/deep0", h); chain143_1(r, p) }
func chain143_1(r Router, p string) { r.Get(p+"/deep1", h); chain143_2(r, p) }
func chain143_2(r Router, p string) { r.Get(p+"/deep2", h); chain143_3(r, p) }
func chain143_3(r Router, p string) { deep_0(r, p) }

type m144 struct{}

func (m144) Register(r Router) { reg144(r, "/m144") }
func reg144(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain144_0(r, p)
}
func chain144_0(r Router, p string) { r.Get(p+"/deep0", h); chain144_1(r, p) }
func chain144_1(r Router, p string) { r.Get(p+"/deep1", h); chain144_2(r, p) }
func chain144_2(r Router, p string) { r.Get(p+"/deep2", h); chain144_3(r, p) }
func chain144_3(r Router, p string) { deep_0(r, p) }

type m145 struct{}

func (m145) Register(r Router) { reg145(r, "/m145") }
func reg145(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain145_0(r, p)
}
func chain145_0(r Router, p string) { r.Get(p+"/deep0", h); chain145_1(r, p) }
func chain145_1(r Router, p string) { r.Get(p+"/deep1", h); chain145_2(r, p) }
func chain145_2(r Router, p string) { r.Get(p+"/deep2", h); chain145_3(r, p) }
func chain145_3(r Router, p string) { deep_0(r, p) }

type m146 struct{}

func (m146) Register(r Router) { reg146(r, "/m146") }
func reg146(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain146_0(r, p)
}
func chain146_0(r Router, p string) { r.Get(p+"/deep0", h); chain146_1(r, p) }
func chain146_1(r Router, p string) { r.Get(p+"/deep1", h); chain146_2(r, p) }
func chain146_2(r Router, p string) { r.Get(p+"/deep2", h); chain146_3(r, p) }
func chain146_3(r Router, p string) { deep_0(r, p) }

type m147 struct{}

func (m147) Register(r Router) { reg147(r, "/m147") }
func reg147(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain147_0(r, p)
}
func chain147_0(r Router, p string) { r.Get(p+"/deep0", h); chain147_1(r, p) }
func chain147_1(r Router, p string) { r.Get(p+"/deep1", h); chain147_2(r, p) }
func chain147_2(r Router, p string) { r.Get(p+"/deep2", h); chain147_3(r, p) }
func chain147_3(r Router, p string) { deep_0(r, p) }

type m148 struct{}

func (m148) Register(r Router) { reg148(r, "/m148") }
func reg148(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain148_0(r, p)
}
func chain148_0(r Router, p string) { r.Get(p+"/deep0", h); chain148_1(r, p) }
func chain148_1(r Router, p string) { r.Get(p+"/deep1", h); chain148_2(r, p) }
func chain148_2(r Router, p string) { r.Get(p+"/deep2", h); chain148_3(r, p) }
func chain148_3(r Router, p string) { deep_0(r, p) }

type m149 struct{}

func (m149) Register(r Router) { reg149(r, "/m149") }
func reg149(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain149_0(r, p)
}
func chain149_0(r Router, p string) { r.Get(p+"/deep0", h); chain149_1(r, p) }
func chain149_1(r Router, p string) { r.Get(p+"/deep1", h); chain149_2(r, p) }
func chain149_2(r Router, p string) { r.Get(p+"/deep2", h); chain149_3(r, p) }
func chain149_3(r Router, p string) { deep_0(r, p) }

type m150 struct{}

func (m150) Register(r Router) { reg150(r, "/m150") }
func reg150(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain150_0(r, p)
}
func chain150_0(r Router, p string) { r.Get(p+"/deep0", h); chain150_1(r, p) }
func chain150_1(r Router, p string) { r.Get(p+"/deep1", h); chain150_2(r, p) }
func chain150_2(r Router, p string) { r.Get(p+"/deep2", h); chain150_3(r, p) }
func chain150_3(r Router, p string) { deep_0(r, p) }

type m151 struct{}

func (m151) Register(r Router) { reg151(r, "/m151") }
func reg151(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain151_0(r, p)
}
func chain151_0(r Router, p string) { r.Get(p+"/deep0", h); chain151_1(r, p) }
func chain151_1(r Router, p string) { r.Get(p+"/deep1", h); chain151_2(r, p) }
func chain151_2(r Router, p string) { r.Get(p+"/deep2", h); chain151_3(r, p) }
func chain151_3(r Router, p string) { deep_0(r, p) }

type m152 struct{}

func (m152) Register(r Router) { reg152(r, "/m152") }
func reg152(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain152_0(r, p)
}
func chain152_0(r Router, p string) { r.Get(p+"/deep0", h); chain152_1(r, p) }
func chain152_1(r Router, p string) { r.Get(p+"/deep1", h); chain152_2(r, p) }
func chain152_2(r Router, p string) { r.Get(p+"/deep2", h); chain152_3(r, p) }
func chain152_3(r Router, p string) { deep_0(r, p) }

type m153 struct{}

func (m153) Register(r Router) { reg153(r, "/m153") }
func reg153(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain153_0(r, p)
}
func chain153_0(r Router, p string) { r.Get(p+"/deep0", h); chain153_1(r, p) }
func chain153_1(r Router, p string) { r.Get(p+"/deep1", h); chain153_2(r, p) }
func chain153_2(r Router, p string) { r.Get(p+"/deep2", h); chain153_3(r, p) }
func chain153_3(r Router, p string) { deep_0(r, p) }

type m154 struct{}

func (m154) Register(r Router) { reg154(r, "/m154") }
func reg154(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain154_0(r, p)
}
func chain154_0(r Router, p string) { r.Get(p+"/deep0", h); chain154_1(r, p) }
func chain154_1(r Router, p string) { r.Get(p+"/deep1", h); chain154_2(r, p) }
func chain154_2(r Router, p string) { r.Get(p+"/deep2", h); chain154_3(r, p) }
func chain154_3(r Router, p string) { deep_0(r, p) }

type m155 struct{}

func (m155) Register(r Router) { reg155(r, "/m155") }
func reg155(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain155_0(r, p)
}
func chain155_0(r Router, p string) { r.Get(p+"/deep0", h); chain155_1(r, p) }
func chain155_1(r Router, p string) { r.Get(p+"/deep1", h); chain155_2(r, p) }
func chain155_2(r Router, p string) { r.Get(p+"/deep2", h); chain155_3(r, p) }
func chain155_3(r Router, p string) { deep_0(r, p) }

type m156 struct{}

func (m156) Register(r Router) { reg156(r, "/m156") }
func reg156(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain156_0(r, p)
}
func chain156_0(r Router, p string) { r.Get(p+"/deep0", h); chain156_1(r, p) }
func chain156_1(r Router, p string) { r.Get(p+"/deep1", h); chain156_2(r, p) }
func chain156_2(r Router, p string) { r.Get(p+"/deep2", h); chain156_3(r, p) }
func chain156_3(r Router, p string) { deep_0(r, p) }

type m157 struct{}

func (m157) Register(r Router) { reg157(r, "/m157") }
func reg157(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain157_0(r, p)
}
func chain157_0(r Router, p string) { r.Get(p+"/deep0", h); chain157_1(r, p) }
func chain157_1(r Router, p string) { r.Get(p+"/deep1", h); chain157_2(r, p) }
func chain157_2(r Router, p string) { r.Get(p+"/deep2", h); chain157_3(r, p) }
func chain157_3(r Router, p string) { deep_0(r, p) }

type m158 struct{}

func (m158) Register(r Router) { reg158(r, "/m158") }
func reg158(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain158_0(r, p)
}
func chain158_0(r Router, p string) { r.Get(p+"/deep0", h); chain158_1(r, p) }
func chain158_1(r Router, p string) { r.Get(p+"/deep1", h); chain158_2(r, p) }
func chain158_2(r Router, p string) { r.Get(p+"/deep2", h); chain158_3(r, p) }
func chain158_3(r Router, p string) { deep_0(r, p) }

type m159 struct{}

func (m159) Register(r Router) { reg159(r, "/m159") }
func reg159(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain159_0(r, p)
}
func chain159_0(r Router, p string) { r.Get(p+"/deep0", h); chain159_1(r, p) }
func chain159_1(r Router, p string) { r.Get(p+"/deep1", h); chain159_2(r, p) }
func chain159_2(r Router, p string) { r.Get(p+"/deep2", h); chain159_3(r, p) }
func chain159_3(r Router, p string) { deep_0(r, p) }

type m160 struct{}

func (m160) Register(r Router) { reg160(r, "/m160") }
func reg160(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain160_0(r, p)
}
func chain160_0(r Router, p string) { r.Get(p+"/deep0", h); chain160_1(r, p) }
func chain160_1(r Router, p string) { r.Get(p+"/deep1", h); chain160_2(r, p) }
func chain160_2(r Router, p string) { r.Get(p+"/deep2", h); chain160_3(r, p) }
func chain160_3(r Router, p string) { deep_0(r, p) }

type m161 struct{}

func (m161) Register(r Router) { reg161(r, "/m161") }
func reg161(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain161_0(r, p)
}
func chain161_0(r Router, p string) { r.Get(p+"/deep0", h); chain161_1(r, p) }
func chain161_1(r Router, p string) { r.Get(p+"/deep1", h); chain161_2(r, p) }
func chain161_2(r Router, p string) { r.Get(p+"/deep2", h); chain161_3(r, p) }
func chain161_3(r Router, p string) { deep_0(r, p) }

type m162 struct{}

func (m162) Register(r Router) { reg162(r, "/m162") }
func reg162(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain162_0(r, p)
}
func chain162_0(r Router, p string) { r.Get(p+"/deep0", h); chain162_1(r, p) }
func chain162_1(r Router, p string) { r.Get(p+"/deep1", h); chain162_2(r, p) }
func chain162_2(r Router, p string) { r.Get(p+"/deep2", h); chain162_3(r, p) }
func chain162_3(r Router, p string) { deep_0(r, p) }

type m163 struct{}

func (m163) Register(r Router) { reg163(r, "/m163") }
func reg163(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain163_0(r, p)
}
func chain163_0(r Router, p string) { r.Get(p+"/deep0", h); chain163_1(r, p) }
func chain163_1(r Router, p string) { r.Get(p+"/deep1", h); chain163_2(r, p) }
func chain163_2(r Router, p string) { r.Get(p+"/deep2", h); chain163_3(r, p) }
func chain163_3(r Router, p string) { deep_0(r, p) }

type m164 struct{}

func (m164) Register(r Router) { reg164(r, "/m164") }
func reg164(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain164_0(r, p)
}
func chain164_0(r Router, p string) { r.Get(p+"/deep0", h); chain164_1(r, p) }
func chain164_1(r Router, p string) { r.Get(p+"/deep1", h); chain164_2(r, p) }
func chain164_2(r Router, p string) { r.Get(p+"/deep2", h); chain164_3(r, p) }
func chain164_3(r Router, p string) { deep_0(r, p) }

type m165 struct{}

func (m165) Register(r Router) { reg165(r, "/m165") }
func reg165(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain165_0(r, p)
}
func chain165_0(r Router, p string) { r.Get(p+"/deep0", h); chain165_1(r, p) }
func chain165_1(r Router, p string) { r.Get(p+"/deep1", h); chain165_2(r, p) }
func chain165_2(r Router, p string) { r.Get(p+"/deep2", h); chain165_3(r, p) }
func chain165_3(r Router, p string) { deep_0(r, p) }

type m166 struct{}

func (m166) Register(r Router) { reg166(r, "/m166") }
func reg166(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain166_0(r, p)
}
func chain166_0(r Router, p string) { r.Get(p+"/deep0", h); chain166_1(r, p) }
func chain166_1(r Router, p string) { r.Get(p+"/deep1", h); chain166_2(r, p) }
func chain166_2(r Router, p string) { r.Get(p+"/deep2", h); chain166_3(r, p) }
func chain166_3(r Router, p string) { deep_0(r, p) }

type m167 struct{}

func (m167) Register(r Router) { reg167(r, "/m167") }
func reg167(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain167_0(r, p)
}
func chain167_0(r Router, p string) { r.Get(p+"/deep0", h); chain167_1(r, p) }
func chain167_1(r Router, p string) { r.Get(p+"/deep1", h); chain167_2(r, p) }
func chain167_2(r Router, p string) { r.Get(p+"/deep2", h); chain167_3(r, p) }
func chain167_3(r Router, p string) { deep_0(r, p) }

type m168 struct{}

func (m168) Register(r Router) { reg168(r, "/m168") }
func reg168(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain168_0(r, p)
}
func chain168_0(r Router, p string) { r.Get(p+"/deep0", h); chain168_1(r, p) }
func chain168_1(r Router, p string) { r.Get(p+"/deep1", h); chain168_2(r, p) }
func chain168_2(r Router, p string) { r.Get(p+"/deep2", h); chain168_3(r, p) }
func chain168_3(r Router, p string) { deep_0(r, p) }

type m169 struct{}

func (m169) Register(r Router) { reg169(r, "/m169") }
func reg169(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain169_0(r, p)
}
func chain169_0(r Router, p string) { r.Get(p+"/deep0", h); chain169_1(r, p) }
func chain169_1(r Router, p string) { r.Get(p+"/deep1", h); chain169_2(r, p) }
func chain169_2(r Router, p string) { r.Get(p+"/deep2", h); chain169_3(r, p) }
func chain169_3(r Router, p string) { deep_0(r, p) }

type m170 struct{}

func (m170) Register(r Router) { reg170(r, "/m170") }
func reg170(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain170_0(r, p)
}
func chain170_0(r Router, p string) { r.Get(p+"/deep0", h); chain170_1(r, p) }
func chain170_1(r Router, p string) { r.Get(p+"/deep1", h); chain170_2(r, p) }
func chain170_2(r Router, p string) { r.Get(p+"/deep2", h); chain170_3(r, p) }
func chain170_3(r Router, p string) { deep_0(r, p) }

type m171 struct{}

func (m171) Register(r Router) { reg171(r, "/m171") }
func reg171(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain171_0(r, p)
}
func chain171_0(r Router, p string) { r.Get(p+"/deep0", h); chain171_1(r, p) }
func chain171_1(r Router, p string) { r.Get(p+"/deep1", h); chain171_2(r, p) }
func chain171_2(r Router, p string) { r.Get(p+"/deep2", h); chain171_3(r, p) }
func chain171_3(r Router, p string) { deep_0(r, p) }

type m172 struct{}

func (m172) Register(r Router) { reg172(r, "/m172") }
func reg172(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain172_0(r, p)
}
func chain172_0(r Router, p string) { r.Get(p+"/deep0", h); chain172_1(r, p) }
func chain172_1(r Router, p string) { r.Get(p+"/deep1", h); chain172_2(r, p) }
func chain172_2(r Router, p string) { r.Get(p+"/deep2", h); chain172_3(r, p) }
func chain172_3(r Router, p string) { deep_0(r, p) }

type m173 struct{}

func (m173) Register(r Router) { reg173(r, "/m173") }
func reg173(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain173_0(r, p)
}
func chain173_0(r Router, p string) { r.Get(p+"/deep0", h); chain173_1(r, p) }
func chain173_1(r Router, p string) { r.Get(p+"/deep1", h); chain173_2(r, p) }
func chain173_2(r Router, p string) { r.Get(p+"/deep2", h); chain173_3(r, p) }
func chain173_3(r Router, p string) { deep_0(r, p) }

type m174 struct{}

func (m174) Register(r Router) { reg174(r, "/m174") }
func reg174(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain174_0(r, p)
}
func chain174_0(r Router, p string) { r.Get(p+"/deep0", h); chain174_1(r, p) }
func chain174_1(r Router, p string) { r.Get(p+"/deep1", h); chain174_2(r, p) }
func chain174_2(r Router, p string) { r.Get(p+"/deep2", h); chain174_3(r, p) }
func chain174_3(r Router, p string) { deep_0(r, p) }

type m175 struct{}

func (m175) Register(r Router) { reg175(r, "/m175") }
func reg175(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain175_0(r, p)
}
func chain175_0(r Router, p string) { r.Get(p+"/deep0", h); chain175_1(r, p) }
func chain175_1(r Router, p string) { r.Get(p+"/deep1", h); chain175_2(r, p) }
func chain175_2(r Router, p string) { r.Get(p+"/deep2", h); chain175_3(r, p) }
func chain175_3(r Router, p string) { deep_0(r, p) }

type m176 struct{}

func (m176) Register(r Router) { reg176(r, "/m176") }
func reg176(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain176_0(r, p)
}
func chain176_0(r Router, p string) { r.Get(p+"/deep0", h); chain176_1(r, p) }
func chain176_1(r Router, p string) { r.Get(p+"/deep1", h); chain176_2(r, p) }
func chain176_2(r Router, p string) { r.Get(p+"/deep2", h); chain176_3(r, p) }
func chain176_3(r Router, p string) { deep_0(r, p) }

type m177 struct{}

func (m177) Register(r Router) { reg177(r, "/m177") }
func reg177(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain177_0(r, p)
}
func chain177_0(r Router, p string) { r.Get(p+"/deep0", h); chain177_1(r, p) }
func chain177_1(r Router, p string) { r.Get(p+"/deep1", h); chain177_2(r, p) }
func chain177_2(r Router, p string) { r.Get(p+"/deep2", h); chain177_3(r, p) }
func chain177_3(r Router, p string) { deep_0(r, p) }

type m178 struct{}

func (m178) Register(r Router) { reg178(r, "/m178") }
func reg178(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain178_0(r, p)
}
func chain178_0(r Router, p string) { r.Get(p+"/deep0", h); chain178_1(r, p) }
func chain178_1(r Router, p string) { r.Get(p+"/deep1", h); chain178_2(r, p) }
func chain178_2(r Router, p string) { r.Get(p+"/deep2", h); chain178_3(r, p) }
func chain178_3(r Router, p string) { deep_0(r, p) }

type m179 struct{}

func (m179) Register(r Router) { reg179(r, "/m179") }
func reg179(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain179_0(r, p)
}
func chain179_0(r Router, p string) { r.Get(p+"/deep0", h); chain179_1(r, p) }
func chain179_1(r Router, p string) { r.Get(p+"/deep1", h); chain179_2(r, p) }
func chain179_2(r Router, p string) { r.Get(p+"/deep2", h); chain179_3(r, p) }
func chain179_3(r Router, p string) { deep_0(r, p) }

type m180 struct{}

func (m180) Register(r Router) { reg180(r, "/m180") }
func reg180(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain180_0(r, p)
}
func chain180_0(r Router, p string) { r.Get(p+"/deep0", h); chain180_1(r, p) }
func chain180_1(r Router, p string) { r.Get(p+"/deep1", h); chain180_2(r, p) }
func chain180_2(r Router, p string) { r.Get(p+"/deep2", h); chain180_3(r, p) }
func chain180_3(r Router, p string) { deep_0(r, p) }

type m181 struct{}

func (m181) Register(r Router) { reg181(r, "/m181") }
func reg181(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain181_0(r, p)
}
func chain181_0(r Router, p string) { r.Get(p+"/deep0", h); chain181_1(r, p) }
func chain181_1(r Router, p string) { r.Get(p+"/deep1", h); chain181_2(r, p) }
func chain181_2(r Router, p string) { r.Get(p+"/deep2", h); chain181_3(r, p) }
func chain181_3(r Router, p string) { deep_0(r, p) }

type m182 struct{}

func (m182) Register(r Router) { reg182(r, "/m182") }
func reg182(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain182_0(r, p)
}
func chain182_0(r Router, p string) { r.Get(p+"/deep0", h); chain182_1(r, p) }
func chain182_1(r Router, p string) { r.Get(p+"/deep1", h); chain182_2(r, p) }
func chain182_2(r Router, p string) { r.Get(p+"/deep2", h); chain182_3(r, p) }
func chain182_3(r Router, p string) { deep_0(r, p) }

type m183 struct{}

func (m183) Register(r Router) { reg183(r, "/m183") }
func reg183(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain183_0(r, p)
}
func chain183_0(r Router, p string) { r.Get(p+"/deep0", h); chain183_1(r, p) }
func chain183_1(r Router, p string) { r.Get(p+"/deep1", h); chain183_2(r, p) }
func chain183_2(r Router, p string) { r.Get(p+"/deep2", h); chain183_3(r, p) }
func chain183_3(r Router, p string) { deep_0(r, p) }

type m184 struct{}

func (m184) Register(r Router) { reg184(r, "/m184") }
func reg184(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain184_0(r, p)
}
func chain184_0(r Router, p string) { r.Get(p+"/deep0", h); chain184_1(r, p) }
func chain184_1(r Router, p string) { r.Get(p+"/deep1", h); chain184_2(r, p) }
func chain184_2(r Router, p string) { r.Get(p+"/deep2", h); chain184_3(r, p) }
func chain184_3(r Router, p string) { deep_0(r, p) }

type m185 struct{}

func (m185) Register(r Router) { reg185(r, "/m185") }
func reg185(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain185_0(r, p)
}
func chain185_0(r Router, p string) { r.Get(p+"/deep0", h); chain185_1(r, p) }
func chain185_1(r Router, p string) { r.Get(p+"/deep1", h); chain185_2(r, p) }
func chain185_2(r Router, p string) { r.Get(p+"/deep2", h); chain185_3(r, p) }
func chain185_3(r Router, p string) { deep_0(r, p) }

type m186 struct{}

func (m186) Register(r Router) { reg186(r, "/m186") }
func reg186(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain186_0(r, p)
}
func chain186_0(r Router, p string) { r.Get(p+"/deep0", h); chain186_1(r, p) }
func chain186_1(r Router, p string) { r.Get(p+"/deep1", h); chain186_2(r, p) }
func chain186_2(r Router, p string) { r.Get(p+"/deep2", h); chain186_3(r, p) }
func chain186_3(r Router, p string) { deep_0(r, p) }

type m187 struct{}

func (m187) Register(r Router) { reg187(r, "/m187") }
func reg187(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain187_0(r, p)
}
func chain187_0(r Router, p string) { r.Get(p+"/deep0", h); chain187_1(r, p) }
func chain187_1(r Router, p string) { r.Get(p+"/deep1", h); chain187_2(r, p) }
func chain187_2(r Router, p string) { r.Get(p+"/deep2", h); chain187_3(r, p) }
func chain187_3(r Router, p string) { deep_0(r, p) }

type m188 struct{}

func (m188) Register(r Router) { reg188(r, "/m188") }
func reg188(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain188_0(r, p)
}
func chain188_0(r Router, p string) { r.Get(p+"/deep0", h); chain188_1(r, p) }
func chain188_1(r Router, p string) { r.Get(p+"/deep1", h); chain188_2(r, p) }
func chain188_2(r Router, p string) { r.Get(p+"/deep2", h); chain188_3(r, p) }
func chain188_3(r Router, p string) { deep_0(r, p) }

type m189 struct{}

func (m189) Register(r Router) { reg189(r, "/m189") }
func reg189(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain189_0(r, p)
}
func chain189_0(r Router, p string) { r.Get(p+"/deep0", h); chain189_1(r, p) }
func chain189_1(r Router, p string) { r.Get(p+"/deep1", h); chain189_2(r, p) }
func chain189_2(r Router, p string) { r.Get(p+"/deep2", h); chain189_3(r, p) }
func chain189_3(r Router, p string) { deep_0(r, p) }

type m190 struct{}

func (m190) Register(r Router) { reg190(r, "/m190") }
func reg190(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain190_0(r, p)
}
func chain190_0(r Router, p string) { r.Get(p+"/deep0", h); chain190_1(r, p) }
func chain190_1(r Router, p string) { r.Get(p+"/deep1", h); chain190_2(r, p) }
func chain190_2(r Router, p string) { r.Get(p+"/deep2", h); chain190_3(r, p) }
func chain190_3(r Router, p string) { deep_0(r, p) }

type m191 struct{}

func (m191) Register(r Router) { reg191(r, "/m191") }
func reg191(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain191_0(r, p)
}
func chain191_0(r Router, p string) { r.Get(p+"/deep0", h); chain191_1(r, p) }
func chain191_1(r Router, p string) { r.Get(p+"/deep1", h); chain191_2(r, p) }
func chain191_2(r Router, p string) { r.Get(p+"/deep2", h); chain191_3(r, p) }
func chain191_3(r Router, p string) { deep_0(r, p) }

type m192 struct{}

func (m192) Register(r Router) { reg192(r, "/m192") }
func reg192(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain192_0(r, p)
}
func chain192_0(r Router, p string) { r.Get(p+"/deep0", h); chain192_1(r, p) }
func chain192_1(r Router, p string) { r.Get(p+"/deep1", h); chain192_2(r, p) }
func chain192_2(r Router, p string) { r.Get(p+"/deep2", h); chain192_3(r, p) }
func chain192_3(r Router, p string) { deep_0(r, p) }

type m193 struct{}

func (m193) Register(r Router) { reg193(r, "/m193") }
func reg193(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain193_0(r, p)
}
func chain193_0(r Router, p string) { r.Get(p+"/deep0", h); chain193_1(r, p) }
func chain193_1(r Router, p string) { r.Get(p+"/deep1", h); chain193_2(r, p) }
func chain193_2(r Router, p string) { r.Get(p+"/deep2", h); chain193_3(r, p) }
func chain193_3(r Router, p string) { deep_0(r, p) }

type m194 struct{}

func (m194) Register(r Router) { reg194(r, "/m194") }
func reg194(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain194_0(r, p)
}
func chain194_0(r Router, p string) { r.Get(p+"/deep0", h); chain194_1(r, p) }
func chain194_1(r Router, p string) { r.Get(p+"/deep1", h); chain194_2(r, p) }
func chain194_2(r Router, p string) { r.Get(p+"/deep2", h); chain194_3(r, p) }
func chain194_3(r Router, p string) { deep_0(r, p) }

type m195 struct{}

func (m195) Register(r Router) { reg195(r, "/m195") }
func reg195(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain195_0(r, p)
}
func chain195_0(r Router, p string) { r.Get(p+"/deep0", h); chain195_1(r, p) }
func chain195_1(r Router, p string) { r.Get(p+"/deep1", h); chain195_2(r, p) }
func chain195_2(r Router, p string) { r.Get(p+"/deep2", h); chain195_3(r, p) }
func chain195_3(r Router, p string) { deep_0(r, p) }

type m196 struct{}

func (m196) Register(r Router) { reg196(r, "/m196") }
func reg196(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain196_0(r, p)
}
func chain196_0(r Router, p string) { r.Get(p+"/deep0", h); chain196_1(r, p) }
func chain196_1(r Router, p string) { r.Get(p+"/deep1", h); chain196_2(r, p) }
func chain196_2(r Router, p string) { r.Get(p+"/deep2", h); chain196_3(r, p) }
func chain196_3(r Router, p string) { deep_0(r, p) }

type m197 struct{}

func (m197) Register(r Router) { reg197(r, "/m197") }
func reg197(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain197_0(r, p)
}
func chain197_0(r Router, p string) { r.Get(p+"/deep0", h); chain197_1(r, p) }
func chain197_1(r Router, p string) { r.Get(p+"/deep1", h); chain197_2(r, p) }
func chain197_2(r Router, p string) { r.Get(p+"/deep2", h); chain197_3(r, p) }
func chain197_3(r Router, p string) { deep_0(r, p) }

type m198 struct{}

func (m198) Register(r Router) { reg198(r, "/m198") }
func reg198(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain198_0(r, p)
}
func chain198_0(r Router, p string) { r.Get(p+"/deep0", h); chain198_1(r, p) }
func chain198_1(r Router, p string) { r.Get(p+"/deep1", h); chain198_2(r, p) }
func chain198_2(r Router, p string) { r.Get(p+"/deep2", h); chain198_3(r, p) }
func chain198_3(r Router, p string) { deep_0(r, p) }

type m199 struct{}

func (m199) Register(r Router) { reg199(r, "/m199") }
func reg199(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain199_0(r, p)
}
func chain199_0(r Router, p string) { r.Get(p+"/deep0", h); chain199_1(r, p) }
func chain199_1(r Router, p string) { r.Get(p+"/deep1", h); chain199_2(r, p) }
func chain199_2(r Router, p string) { r.Get(p+"/deep2", h); chain199_3(r, p) }
func chain199_3(r Router, p string) { deep_0(r, p) }

type m200 struct{}

func (m200) Register(r Router) { reg200(r, "/m200") }
func reg200(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain200_0(r, p)
}
func chain200_0(r Router, p string) { r.Get(p+"/deep0", h); chain200_1(r, p) }
func chain200_1(r Router, p string) { r.Get(p+"/deep1", h); chain200_2(r, p) }
func chain200_2(r Router, p string) { r.Get(p+"/deep2", h); chain200_3(r, p) }
func chain200_3(r Router, p string) { deep_0(r, p) }

type m201 struct{}

func (m201) Register(r Router) { reg201(r, "/m201") }
func reg201(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain201_0(r, p)
}
func chain201_0(r Router, p string) { r.Get(p+"/deep0", h); chain201_1(r, p) }
func chain201_1(r Router, p string) { r.Get(p+"/deep1", h); chain201_2(r, p) }
func chain201_2(r Router, p string) { r.Get(p+"/deep2", h); chain201_3(r, p) }
func chain201_3(r Router, p string) { deep_0(r, p) }

type m202 struct{}

func (m202) Register(r Router) { reg202(r, "/m202") }
func reg202(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain202_0(r, p)
}
func chain202_0(r Router, p string) { r.Get(p+"/deep0", h); chain202_1(r, p) }
func chain202_1(r Router, p string) { r.Get(p+"/deep1", h); chain202_2(r, p) }
func chain202_2(r Router, p string) { r.Get(p+"/deep2", h); chain202_3(r, p) }
func chain202_3(r Router, p string) { deep_0(r, p) }

type m203 struct{}

func (m203) Register(r Router) { reg203(r, "/m203") }
func reg203(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain203_0(r, p)
}
func chain203_0(r Router, p string) { r.Get(p+"/deep0", h); chain203_1(r, p) }
func chain203_1(r Router, p string) { r.Get(p+"/deep1", h); chain203_2(r, p) }
func chain203_2(r Router, p string) { r.Get(p+"/deep2", h); chain203_3(r, p) }
func chain203_3(r Router, p string) { deep_0(r, p) }

type m204 struct{}

func (m204) Register(r Router) { reg204(r, "/m204") }
func reg204(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain204_0(r, p)
}
func chain204_0(r Router, p string) { r.Get(p+"/deep0", h); chain204_1(r, p) }
func chain204_1(r Router, p string) { r.Get(p+"/deep1", h); chain204_2(r, p) }
func chain204_2(r Router, p string) { r.Get(p+"/deep2", h); chain204_3(r, p) }
func chain204_3(r Router, p string) { deep_0(r, p) }

type m205 struct{}

func (m205) Register(r Router) { reg205(r, "/m205") }
func reg205(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain205_0(r, p)
}
func chain205_0(r Router, p string) { r.Get(p+"/deep0", h); chain205_1(r, p) }
func chain205_1(r Router, p string) { r.Get(p+"/deep1", h); chain205_2(r, p) }
func chain205_2(r Router, p string) { r.Get(p+"/deep2", h); chain205_3(r, p) }
func chain205_3(r Router, p string) { deep_0(r, p) }

type m206 struct{}

func (m206) Register(r Router) { reg206(r, "/m206") }
func reg206(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain206_0(r, p)
}
func chain206_0(r Router, p string) { r.Get(p+"/deep0", h); chain206_1(r, p) }
func chain206_1(r Router, p string) { r.Get(p+"/deep1", h); chain206_2(r, p) }
func chain206_2(r Router, p string) { r.Get(p+"/deep2", h); chain206_3(r, p) }
func chain206_3(r Router, p string) { deep_0(r, p) }

type m207 struct{}

func (m207) Register(r Router) { reg207(r, "/m207") }
func reg207(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain207_0(r, p)
}
func chain207_0(r Router, p string) { r.Get(p+"/deep0", h); chain207_1(r, p) }
func chain207_1(r Router, p string) { r.Get(p+"/deep1", h); chain207_2(r, p) }
func chain207_2(r Router, p string) { r.Get(p+"/deep2", h); chain207_3(r, p) }
func chain207_3(r Router, p string) { deep_0(r, p) }

type m208 struct{}

func (m208) Register(r Router) { reg208(r, "/m208") }
func reg208(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain208_0(r, p)
}
func chain208_0(r Router, p string) { r.Get(p+"/deep0", h); chain208_1(r, p) }
func chain208_1(r Router, p string) { r.Get(p+"/deep1", h); chain208_2(r, p) }
func chain208_2(r Router, p string) { r.Get(p+"/deep2", h); chain208_3(r, p) }
func chain208_3(r Router, p string) { deep_0(r, p) }

type m209 struct{}

func (m209) Register(r Router) { reg209(r, "/m209") }
func reg209(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain209_0(r, p)
}
func chain209_0(r Router, p string) { r.Get(p+"/deep0", h); chain209_1(r, p) }
func chain209_1(r Router, p string) { r.Get(p+"/deep1", h); chain209_2(r, p) }
func chain209_2(r Router, p string) { r.Get(p+"/deep2", h); chain209_3(r, p) }
func chain209_3(r Router, p string) { deep_0(r, p) }

type m210 struct{}

func (m210) Register(r Router) { reg210(r, "/m210") }
func reg210(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain210_0(r, p)
}
func chain210_0(r Router, p string) { r.Get(p+"/deep0", h); chain210_1(r, p) }
func chain210_1(r Router, p string) { r.Get(p+"/deep1", h); chain210_2(r, p) }
func chain210_2(r Router, p string) { r.Get(p+"/deep2", h); chain210_3(r, p) }
func chain210_3(r Router, p string) { deep_0(r, p) }

type m211 struct{}

func (m211) Register(r Router) { reg211(r, "/m211") }
func reg211(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain211_0(r, p)
}
func chain211_0(r Router, p string) { r.Get(p+"/deep0", h); chain211_1(r, p) }
func chain211_1(r Router, p string) { r.Get(p+"/deep1", h); chain211_2(r, p) }
func chain211_2(r Router, p string) { r.Get(p+"/deep2", h); chain211_3(r, p) }
func chain211_3(r Router, p string) { deep_0(r, p) }

type m212 struct{}

func (m212) Register(r Router) { reg212(r, "/m212") }
func reg212(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain212_0(r, p)
}
func chain212_0(r Router, p string) { r.Get(p+"/deep0", h); chain212_1(r, p) }
func chain212_1(r Router, p string) { r.Get(p+"/deep1", h); chain212_2(r, p) }
func chain212_2(r Router, p string) { r.Get(p+"/deep2", h); chain212_3(r, p) }
func chain212_3(r Router, p string) { deep_0(r, p) }

type m213 struct{}

func (m213) Register(r Router) { reg213(r, "/m213") }
func reg213(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain213_0(r, p)
}
func chain213_0(r Router, p string) { r.Get(p+"/deep0", h); chain213_1(r, p) }
func chain213_1(r Router, p string) { r.Get(p+"/deep1", h); chain213_2(r, p) }
func chain213_2(r Router, p string) { r.Get(p+"/deep2", h); chain213_3(r, p) }
func chain213_3(r Router, p string) { deep_0(r, p) }

type m214 struct{}

func (m214) Register(r Router) { reg214(r, "/m214") }
func reg214(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain214_0(r, p)
}
func chain214_0(r Router, p string) { r.Get(p+"/deep0", h); chain214_1(r, p) }
func chain214_1(r Router, p string) { r.Get(p+"/deep1", h); chain214_2(r, p) }
func chain214_2(r Router, p string) { r.Get(p+"/deep2", h); chain214_3(r, p) }
func chain214_3(r Router, p string) { deep_0(r, p) }

type m215 struct{}

func (m215) Register(r Router) { reg215(r, "/m215") }
func reg215(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain215_0(r, p)
}
func chain215_0(r Router, p string) { r.Get(p+"/deep0", h); chain215_1(r, p) }
func chain215_1(r Router, p string) { r.Get(p+"/deep1", h); chain215_2(r, p) }
func chain215_2(r Router, p string) { r.Get(p+"/deep2", h); chain215_3(r, p) }
func chain215_3(r Router, p string) { deep_0(r, p) }

type m216 struct{}

func (m216) Register(r Router) { reg216(r, "/m216") }
func reg216(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain216_0(r, p)
}
func chain216_0(r Router, p string) { r.Get(p+"/deep0", h); chain216_1(r, p) }
func chain216_1(r Router, p string) { r.Get(p+"/deep1", h); chain216_2(r, p) }
func chain216_2(r Router, p string) { r.Get(p+"/deep2", h); chain216_3(r, p) }
func chain216_3(r Router, p string) { deep_0(r, p) }

type m217 struct{}

func (m217) Register(r Router) { reg217(r, "/m217") }
func reg217(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain217_0(r, p)
}
func chain217_0(r Router, p string) { r.Get(p+"/deep0", h); chain217_1(r, p) }
func chain217_1(r Router, p string) { r.Get(p+"/deep1", h); chain217_2(r, p) }
func chain217_2(r Router, p string) { r.Get(p+"/deep2", h); chain217_3(r, p) }
func chain217_3(r Router, p string) { deep_0(r, p) }

type m218 struct{}

func (m218) Register(r Router) { reg218(r, "/m218") }
func reg218(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain218_0(r, p)
}
func chain218_0(r Router, p string) { r.Get(p+"/deep0", h); chain218_1(r, p) }
func chain218_1(r Router, p string) { r.Get(p+"/deep1", h); chain218_2(r, p) }
func chain218_2(r Router, p string) { r.Get(p+"/deep2", h); chain218_3(r, p) }
func chain218_3(r Router, p string) { deep_0(r, p) }

type m219 struct{}

func (m219) Register(r Router) { reg219(r, "/m219") }
func reg219(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain219_0(r, p)
}
func chain219_0(r Router, p string) { r.Get(p+"/deep0", h); chain219_1(r, p) }
func chain219_1(r Router, p string) { r.Get(p+"/deep1", h); chain219_2(r, p) }
func chain219_2(r Router, p string) { r.Get(p+"/deep2", h); chain219_3(r, p) }
func chain219_3(r Router, p string) { deep_0(r, p) }

type m220 struct{}

func (m220) Register(r Router) { reg220(r, "/m220") }
func reg220(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain220_0(r, p)
}
func chain220_0(r Router, p string) { r.Get(p+"/deep0", h); chain220_1(r, p) }
func chain220_1(r Router, p string) { r.Get(p+"/deep1", h); chain220_2(r, p) }
func chain220_2(r Router, p string) { r.Get(p+"/deep2", h); chain220_3(r, p) }
func chain220_3(r Router, p string) { deep_0(r, p) }

type m221 struct{}

func (m221) Register(r Router) { reg221(r, "/m221") }
func reg221(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain221_0(r, p)
}
func chain221_0(r Router, p string) { r.Get(p+"/deep0", h); chain221_1(r, p) }
func chain221_1(r Router, p string) { r.Get(p+"/deep1", h); chain221_2(r, p) }
func chain221_2(r Router, p string) { r.Get(p+"/deep2", h); chain221_3(r, p) }
func chain221_3(r Router, p string) { deep_0(r, p) }

type m222 struct{}

func (m222) Register(r Router) { reg222(r, "/m222") }
func reg222(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain222_0(r, p)
}
func chain222_0(r Router, p string) { r.Get(p+"/deep0", h); chain222_1(r, p) }
func chain222_1(r Router, p string) { r.Get(p+"/deep1", h); chain222_2(r, p) }
func chain222_2(r Router, p string) { r.Get(p+"/deep2", h); chain222_3(r, p) }
func chain222_3(r Router, p string) { deep_0(r, p) }

type m223 struct{}

func (m223) Register(r Router) { reg223(r, "/m223") }
func reg223(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain223_0(r, p)
}
func chain223_0(r Router, p string) { r.Get(p+"/deep0", h); chain223_1(r, p) }
func chain223_1(r Router, p string) { r.Get(p+"/deep1", h); chain223_2(r, p) }
func chain223_2(r Router, p string) { r.Get(p+"/deep2", h); chain223_3(r, p) }
func chain223_3(r Router, p string) { deep_0(r, p) }

type m224 struct{}

func (m224) Register(r Router) { reg224(r, "/m224") }
func reg224(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain224_0(r, p)
}
func chain224_0(r Router, p string) { r.Get(p+"/deep0", h); chain224_1(r, p) }
func chain224_1(r Router, p string) { r.Get(p+"/deep1", h); chain224_2(r, p) }
func chain224_2(r Router, p string) { r.Get(p+"/deep2", h); chain224_3(r, p) }
func chain224_3(r Router, p string) { deep_0(r, p) }

type m225 struct{}

func (m225) Register(r Router) { reg225(r, "/m225") }
func reg225(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain225_0(r, p)
}
func chain225_0(r Router, p string) { r.Get(p+"/deep0", h); chain225_1(r, p) }
func chain225_1(r Router, p string) { r.Get(p+"/deep1", h); chain225_2(r, p) }
func chain225_2(r Router, p string) { r.Get(p+"/deep2", h); chain225_3(r, p) }
func chain225_3(r Router, p string) { deep_0(r, p) }

type m226 struct{}

func (m226) Register(r Router) { reg226(r, "/m226") }
func reg226(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain226_0(r, p)
}
func chain226_0(r Router, p string) { r.Get(p+"/deep0", h); chain226_1(r, p) }
func chain226_1(r Router, p string) { r.Get(p+"/deep1", h); chain226_2(r, p) }
func chain226_2(r Router, p string) { r.Get(p+"/deep2", h); chain226_3(r, p) }
func chain226_3(r Router, p string) { deep_0(r, p) }

type m227 struct{}

func (m227) Register(r Router) { reg227(r, "/m227") }
func reg227(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain227_0(r, p)
}
func chain227_0(r Router, p string) { r.Get(p+"/deep0", h); chain227_1(r, p) }
func chain227_1(r Router, p string) { r.Get(p+"/deep1", h); chain227_2(r, p) }
func chain227_2(r Router, p string) { r.Get(p+"/deep2", h); chain227_3(r, p) }
func chain227_3(r Router, p string) { deep_0(r, p) }

type m228 struct{}

func (m228) Register(r Router) { reg228(r, "/m228") }
func reg228(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain228_0(r, p)
}
func chain228_0(r Router, p string) { r.Get(p+"/deep0", h); chain228_1(r, p) }
func chain228_1(r Router, p string) { r.Get(p+"/deep1", h); chain228_2(r, p) }
func chain228_2(r Router, p string) { r.Get(p+"/deep2", h); chain228_3(r, p) }
func chain228_3(r Router, p string) { deep_0(r, p) }

type m229 struct{}

func (m229) Register(r Router) { reg229(r, "/m229") }
func reg229(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain229_0(r, p)
}
func chain229_0(r Router, p string) { r.Get(p+"/deep0", h); chain229_1(r, p) }
func chain229_1(r Router, p string) { r.Get(p+"/deep1", h); chain229_2(r, p) }
func chain229_2(r Router, p string) { r.Get(p+"/deep2", h); chain229_3(r, p) }
func chain229_3(r Router, p string) { deep_0(r, p) }

type m230 struct{}

func (m230) Register(r Router) { reg230(r, "/m230") }
func reg230(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain230_0(r, p)
}
func chain230_0(r Router, p string) { r.Get(p+"/deep0", h); chain230_1(r, p) }
func chain230_1(r Router, p string) { r.Get(p+"/deep1", h); chain230_2(r, p) }
func chain230_2(r Router, p string) { r.Get(p+"/deep2", h); chain230_3(r, p) }
func chain230_3(r Router, p string) { deep_0(r, p) }

type m231 struct{}

func (m231) Register(r Router) { reg231(r, "/m231") }
func reg231(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain231_0(r, p)
}
func chain231_0(r Router, p string) { r.Get(p+"/deep0", h); chain231_1(r, p) }
func chain231_1(r Router, p string) { r.Get(p+"/deep1", h); chain231_2(r, p) }
func chain231_2(r Router, p string) { r.Get(p+"/deep2", h); chain231_3(r, p) }
func chain231_3(r Router, p string) { deep_0(r, p) }

type m232 struct{}

func (m232) Register(r Router) { reg232(r, "/m232") }
func reg232(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain232_0(r, p)
}
func chain232_0(r Router, p string) { r.Get(p+"/deep0", h); chain232_1(r, p) }
func chain232_1(r Router, p string) { r.Get(p+"/deep1", h); chain232_2(r, p) }
func chain232_2(r Router, p string) { r.Get(p+"/deep2", h); chain232_3(r, p) }
func chain232_3(r Router, p string) { deep_0(r, p) }

type m233 struct{}

func (m233) Register(r Router) { reg233(r, "/m233") }
func reg233(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain233_0(r, p)
}
func chain233_0(r Router, p string) { r.Get(p+"/deep0", h); chain233_1(r, p) }
func chain233_1(r Router, p string) { r.Get(p+"/deep1", h); chain233_2(r, p) }
func chain233_2(r Router, p string) { r.Get(p+"/deep2", h); chain233_3(r, p) }
func chain233_3(r Router, p string) { deep_0(r, p) }

type m234 struct{}

func (m234) Register(r Router) { reg234(r, "/m234") }
func reg234(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain234_0(r, p)
}
func chain234_0(r Router, p string) { r.Get(p+"/deep0", h); chain234_1(r, p) }
func chain234_1(r Router, p string) { r.Get(p+"/deep1", h); chain234_2(r, p) }
func chain234_2(r Router, p string) { r.Get(p+"/deep2", h); chain234_3(r, p) }
func chain234_3(r Router, p string) { deep_0(r, p) }

type m235 struct{}

func (m235) Register(r Router) { reg235(r, "/m235") }
func reg235(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain235_0(r, p)
}
func chain235_0(r Router, p string) { r.Get(p+"/deep0", h); chain235_1(r, p) }
func chain235_1(r Router, p string) { r.Get(p+"/deep1", h); chain235_2(r, p) }
func chain235_2(r Router, p string) { r.Get(p+"/deep2", h); chain235_3(r, p) }
func chain235_3(r Router, p string) { deep_0(r, p) }

type m236 struct{}

func (m236) Register(r Router) { reg236(r, "/m236") }
func reg236(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain236_0(r, p)
}
func chain236_0(r Router, p string) { r.Get(p+"/deep0", h); chain236_1(r, p) }
func chain236_1(r Router, p string) { r.Get(p+"/deep1", h); chain236_2(r, p) }
func chain236_2(r Router, p string) { r.Get(p+"/deep2", h); chain236_3(r, p) }
func chain236_3(r Router, p string) { deep_0(r, p) }

type m237 struct{}

func (m237) Register(r Router) { reg237(r, "/m237") }
func reg237(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain237_0(r, p)
}
func chain237_0(r Router, p string) { r.Get(p+"/deep0", h); chain237_1(r, p) }
func chain237_1(r Router, p string) { r.Get(p+"/deep1", h); chain237_2(r, p) }
func chain237_2(r Router, p string) { r.Get(p+"/deep2", h); chain237_3(r, p) }
func chain237_3(r Router, p string) { deep_0(r, p) }

type m238 struct{}

func (m238) Register(r Router) { reg238(r, "/m238") }
func reg238(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain238_0(r, p)
}
func chain238_0(r Router, p string) { r.Get(p+"/deep0", h); chain238_1(r, p) }
func chain238_1(r Router, p string) { r.Get(p+"/deep1", h); chain238_2(r, p) }
func chain238_2(r Router, p string) { r.Get(p+"/deep2", h); chain238_3(r, p) }
func chain238_3(r Router, p string) { deep_0(r, p) }

type m239 struct{}

func (m239) Register(r Router) { reg239(r, "/m239") }
func reg239(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain239_0(r, p)
}
func chain239_0(r Router, p string) { r.Get(p+"/deep0", h); chain239_1(r, p) }
func chain239_1(r Router, p string) { r.Get(p+"/deep1", h); chain239_2(r, p) }
func chain239_2(r Router, p string) { r.Get(p+"/deep2", h); chain239_3(r, p) }
func chain239_3(r Router, p string) { deep_0(r, p) }

type m240 struct{}

func (m240) Register(r Router) { reg240(r, "/m240") }
func reg240(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain240_0(r, p)
}
func chain240_0(r Router, p string) { r.Get(p+"/deep0", h); chain240_1(r, p) }
func chain240_1(r Router, p string) { r.Get(p+"/deep1", h); chain240_2(r, p) }
func chain240_2(r Router, p string) { r.Get(p+"/deep2", h); chain240_3(r, p) }
func chain240_3(r Router, p string) { deep_0(r, p) }

type m241 struct{}

func (m241) Register(r Router) { reg241(r, "/m241") }
func reg241(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain241_0(r, p)
}
func chain241_0(r Router, p string) { r.Get(p+"/deep0", h); chain241_1(r, p) }
func chain241_1(r Router, p string) { r.Get(p+"/deep1", h); chain241_2(r, p) }
func chain241_2(r Router, p string) { r.Get(p+"/deep2", h); chain241_3(r, p) }
func chain241_3(r Router, p string) { deep_0(r, p) }

type m242 struct{}

func (m242) Register(r Router) { reg242(r, "/m242") }
func reg242(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain242_0(r, p)
}
func chain242_0(r Router, p string) { r.Get(p+"/deep0", h); chain242_1(r, p) }
func chain242_1(r Router, p string) { r.Get(p+"/deep1", h); chain242_2(r, p) }
func chain242_2(r Router, p string) { r.Get(p+"/deep2", h); chain242_3(r, p) }
func chain242_3(r Router, p string) { deep_0(r, p) }

type m243 struct{}

func (m243) Register(r Router) { reg243(r, "/m243") }
func reg243(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain243_0(r, p)
}
func chain243_0(r Router, p string) { r.Get(p+"/deep0", h); chain243_1(r, p) }
func chain243_1(r Router, p string) { r.Get(p+"/deep1", h); chain243_2(r, p) }
func chain243_2(r Router, p string) { r.Get(p+"/deep2", h); chain243_3(r, p) }
func chain243_3(r Router, p string) { deep_0(r, p) }

type m244 struct{}

func (m244) Register(r Router) { reg244(r, "/m244") }
func reg244(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain244_0(r, p)
}
func chain244_0(r Router, p string) { r.Get(p+"/deep0", h); chain244_1(r, p) }
func chain244_1(r Router, p string) { r.Get(p+"/deep1", h); chain244_2(r, p) }
func chain244_2(r Router, p string) { r.Get(p+"/deep2", h); chain244_3(r, p) }
func chain244_3(r Router, p string) { deep_0(r, p) }

type m245 struct{}

func (m245) Register(r Router) { reg245(r, "/m245") }
func reg245(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain245_0(r, p)
}
func chain245_0(r Router, p string) { r.Get(p+"/deep0", h); chain245_1(r, p) }
func chain245_1(r Router, p string) { r.Get(p+"/deep1", h); chain245_2(r, p) }
func chain245_2(r Router, p string) { r.Get(p+"/deep2", h); chain245_3(r, p) }
func chain245_3(r Router, p string) { deep_0(r, p) }

type m246 struct{}

func (m246) Register(r Router) { reg246(r, "/m246") }
func reg246(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain246_0(r, p)
}
func chain246_0(r Router, p string) { r.Get(p+"/deep0", h); chain246_1(r, p) }
func chain246_1(r Router, p string) { r.Get(p+"/deep1", h); chain246_2(r, p) }
func chain246_2(r Router, p string) { r.Get(p+"/deep2", h); chain246_3(r, p) }
func chain246_3(r Router, p string) { deep_0(r, p) }

type m247 struct{}

func (m247) Register(r Router) { reg247(r, "/m247") }
func reg247(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain247_0(r, p)
}
func chain247_0(r Router, p string) { r.Get(p+"/deep0", h); chain247_1(r, p) }
func chain247_1(r Router, p string) { r.Get(p+"/deep1", h); chain247_2(r, p) }
func chain247_2(r Router, p string) { r.Get(p+"/deep2", h); chain247_3(r, p) }
func chain247_3(r Router, p string) { deep_0(r, p) }

type m248 struct{}

func (m248) Register(r Router) { reg248(r, "/m248") }
func reg248(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain248_0(r, p)
}
func chain248_0(r Router, p string) { r.Get(p+"/deep0", h); chain248_1(r, p) }
func chain248_1(r Router, p string) { r.Get(p+"/deep1", h); chain248_2(r, p) }
func chain248_2(r Router, p string) { r.Get(p+"/deep2", h); chain248_3(r, p) }
func chain248_3(r Router, p string) { deep_0(r, p) }

type m249 struct{}

func (m249) Register(r Router) { reg249(r, "/m249") }
func reg249(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain249_0(r, p)
}
func chain249_0(r Router, p string) { r.Get(p+"/deep0", h); chain249_1(r, p) }
func chain249_1(r Router, p string) { r.Get(p+"/deep1", h); chain249_2(r, p) }
func chain249_2(r Router, p string) { r.Get(p+"/deep2", h); chain249_3(r, p) }
func chain249_3(r Router, p string) { deep_0(r, p) }

type m250 struct{}

func (m250) Register(r Router) { reg250(r, "/m250") }
func reg250(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain250_0(r, p)
}
func chain250_0(r Router, p string) { r.Get(p+"/deep0", h); chain250_1(r, p) }
func chain250_1(r Router, p string) { r.Get(p+"/deep1", h); chain250_2(r, p) }
func chain250_2(r Router, p string) { r.Get(p+"/deep2", h); chain250_3(r, p) }
func chain250_3(r Router, p string) { deep_0(r, p) }

type m251 struct{}

func (m251) Register(r Router) { reg251(r, "/m251") }
func reg251(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain251_0(r, p)
}
func chain251_0(r Router, p string) { r.Get(p+"/deep0", h); chain251_1(r, p) }
func chain251_1(r Router, p string) { r.Get(p+"/deep1", h); chain251_2(r, p) }
func chain251_2(r Router, p string) { r.Get(p+"/deep2", h); chain251_3(r, p) }
func chain251_3(r Router, p string) { deep_0(r, p) }

type m252 struct{}

func (m252) Register(r Router) { reg252(r, "/m252") }
func reg252(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain252_0(r, p)
}
func chain252_0(r Router, p string) { r.Get(p+"/deep0", h); chain252_1(r, p) }
func chain252_1(r Router, p string) { r.Get(p+"/deep1", h); chain252_2(r, p) }
func chain252_2(r Router, p string) { r.Get(p+"/deep2", h); chain252_3(r, p) }
func chain252_3(r Router, p string) { deep_0(r, p) }

type m253 struct{}

func (m253) Register(r Router) { reg253(r, "/m253") }
func reg253(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain253_0(r, p)
}
func chain253_0(r Router, p string) { r.Get(p+"/deep0", h); chain253_1(r, p) }
func chain253_1(r Router, p string) { r.Get(p+"/deep1", h); chain253_2(r, p) }
func chain253_2(r Router, p string) { r.Get(p+"/deep2", h); chain253_3(r, p) }
func chain253_3(r Router, p string) { deep_0(r, p) }

type m254 struct{}

func (m254) Register(r Router) { reg254(r, "/m254") }
func reg254(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain254_0(r, p)
}
func chain254_0(r Router, p string) { r.Get(p+"/deep0", h); chain254_1(r, p) }
func chain254_1(r Router, p string) { r.Get(p+"/deep1", h); chain254_2(r, p) }
func chain254_2(r Router, p string) { r.Get(p+"/deep2", h); chain254_3(r, p) }
func chain254_3(r Router, p string) { deep_0(r, p) }

type m255 struct{}

func (m255) Register(r Router) { reg255(r, "/m255") }
func reg255(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain255_0(r, p)
}
func chain255_0(r Router, p string) { r.Get(p+"/deep0", h); chain255_1(r, p) }
func chain255_1(r Router, p string) { r.Get(p+"/deep1", h); chain255_2(r, p) }
func chain255_2(r Router, p string) { r.Get(p+"/deep2", h); chain255_3(r, p) }
func chain255_3(r Router, p string) { deep_0(r, p) }

type m256 struct{}

func (m256) Register(r Router) { reg256(r, "/m256") }
func reg256(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain256_0(r, p)
}
func chain256_0(r Router, p string) { r.Get(p+"/deep0", h); chain256_1(r, p) }
func chain256_1(r Router, p string) { r.Get(p+"/deep1", h); chain256_2(r, p) }
func chain256_2(r Router, p string) { r.Get(p+"/deep2", h); chain256_3(r, p) }
func chain256_3(r Router, p string) { deep_0(r, p) }

type m257 struct{}

func (m257) Register(r Router) { reg257(r, "/m257") }
func reg257(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain257_0(r, p)
}
func chain257_0(r Router, p string) { r.Get(p+"/deep0", h); chain257_1(r, p) }
func chain257_1(r Router, p string) { r.Get(p+"/deep1", h); chain257_2(r, p) }
func chain257_2(r Router, p string) { r.Get(p+"/deep2", h); chain257_3(r, p) }
func chain257_3(r Router, p string) { deep_0(r, p) }

type m258 struct{}

func (m258) Register(r Router) { reg258(r, "/m258") }
func reg258(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain258_0(r, p)
}
func chain258_0(r Router, p string) { r.Get(p+"/deep0", h); chain258_1(r, p) }
func chain258_1(r Router, p string) { r.Get(p+"/deep1", h); chain258_2(r, p) }
func chain258_2(r Router, p string) { r.Get(p+"/deep2", h); chain258_3(r, p) }
func chain258_3(r Router, p string) { deep_0(r, p) }

type m259 struct{}

func (m259) Register(r Router) { reg259(r, "/m259") }
func reg259(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain259_0(r, p)
}
func chain259_0(r Router, p string) { r.Get(p+"/deep0", h); chain259_1(r, p) }
func chain259_1(r Router, p string) { r.Get(p+"/deep1", h); chain259_2(r, p) }
func chain259_2(r Router, p string) { r.Get(p+"/deep2", h); chain259_3(r, p) }
func chain259_3(r Router, p string) { deep_0(r, p) }

type m260 struct{}

func (m260) Register(r Router) { reg260(r, "/m260") }
func reg260(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain260_0(r, p)
}
func chain260_0(r Router, p string) { r.Get(p+"/deep0", h); chain260_1(r, p) }
func chain260_1(r Router, p string) { r.Get(p+"/deep1", h); chain260_2(r, p) }
func chain260_2(r Router, p string) { r.Get(p+"/deep2", h); chain260_3(r, p) }
func chain260_3(r Router, p string) { deep_0(r, p) }

type m261 struct{}

func (m261) Register(r Router) { reg261(r, "/m261") }
func reg261(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain261_0(r, p)
}
func chain261_0(r Router, p string) { r.Get(p+"/deep0", h); chain261_1(r, p) }
func chain261_1(r Router, p string) { r.Get(p+"/deep1", h); chain261_2(r, p) }
func chain261_2(r Router, p string) { r.Get(p+"/deep2", h); chain261_3(r, p) }
func chain261_3(r Router, p string) { deep_0(r, p) }

type m262 struct{}

func (m262) Register(r Router) { reg262(r, "/m262") }
func reg262(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain262_0(r, p)
}
func chain262_0(r Router, p string) { r.Get(p+"/deep0", h); chain262_1(r, p) }
func chain262_1(r Router, p string) { r.Get(p+"/deep1", h); chain262_2(r, p) }
func chain262_2(r Router, p string) { r.Get(p+"/deep2", h); chain262_3(r, p) }
func chain262_3(r Router, p string) { deep_0(r, p) }

type m263 struct{}

func (m263) Register(r Router) { reg263(r, "/m263") }
func reg263(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain263_0(r, p)
}
func chain263_0(r Router, p string) { r.Get(p+"/deep0", h); chain263_1(r, p) }
func chain263_1(r Router, p string) { r.Get(p+"/deep1", h); chain263_2(r, p) }
func chain263_2(r Router, p string) { r.Get(p+"/deep2", h); chain263_3(r, p) }
func chain263_3(r Router, p string) { deep_0(r, p) }

type m264 struct{}

func (m264) Register(r Router) { reg264(r, "/m264") }
func reg264(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain264_0(r, p)
}
func chain264_0(r Router, p string) { r.Get(p+"/deep0", h); chain264_1(r, p) }
func chain264_1(r Router, p string) { r.Get(p+"/deep1", h); chain264_2(r, p) }
func chain264_2(r Router, p string) { r.Get(p+"/deep2", h); chain264_3(r, p) }
func chain264_3(r Router, p string) { deep_0(r, p) }

type m265 struct{}

func (m265) Register(r Router) { reg265(r, "/m265") }
func reg265(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain265_0(r, p)
}
func chain265_0(r Router, p string) { r.Get(p+"/deep0", h); chain265_1(r, p) }
func chain265_1(r Router, p string) { r.Get(p+"/deep1", h); chain265_2(r, p) }
func chain265_2(r Router, p string) { r.Get(p+"/deep2", h); chain265_3(r, p) }
func chain265_3(r Router, p string) { deep_0(r, p) }

type m266 struct{}

func (m266) Register(r Router) { reg266(r, "/m266") }
func reg266(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain266_0(r, p)
}
func chain266_0(r Router, p string) { r.Get(p+"/deep0", h); chain266_1(r, p) }
func chain266_1(r Router, p string) { r.Get(p+"/deep1", h); chain266_2(r, p) }
func chain266_2(r Router, p string) { r.Get(p+"/deep2", h); chain266_3(r, p) }
func chain266_3(r Router, p string) { deep_0(r, p) }

type m267 struct{}

func (m267) Register(r Router) { reg267(r, "/m267") }
func reg267(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain267_0(r, p)
}
func chain267_0(r Router, p string) { r.Get(p+"/deep0", h); chain267_1(r, p) }
func chain267_1(r Router, p string) { r.Get(p+"/deep1", h); chain267_2(r, p) }
func chain267_2(r Router, p string) { r.Get(p+"/deep2", h); chain267_3(r, p) }
func chain267_3(r Router, p string) { deep_0(r, p) }

type m268 struct{}

func (m268) Register(r Router) { reg268(r, "/m268") }
func reg268(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain268_0(r, p)
}
func chain268_0(r Router, p string) { r.Get(p+"/deep0", h); chain268_1(r, p) }
func chain268_1(r Router, p string) { r.Get(p+"/deep1", h); chain268_2(r, p) }
func chain268_2(r Router, p string) { r.Get(p+"/deep2", h); chain268_3(r, p) }
func chain268_3(r Router, p string) { deep_0(r, p) }

type m269 struct{}

func (m269) Register(r Router) { reg269(r, "/m269") }
func reg269(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain269_0(r, p)
}
func chain269_0(r Router, p string) { r.Get(p+"/deep0", h); chain269_1(r, p) }
func chain269_1(r Router, p string) { r.Get(p+"/deep1", h); chain269_2(r, p) }
func chain269_2(r Router, p string) { r.Get(p+"/deep2", h); chain269_3(r, p) }
func chain269_3(r Router, p string) { deep_0(r, p) }

type m270 struct{}

func (m270) Register(r Router) { reg270(r, "/m270") }
func reg270(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain270_0(r, p)
}
func chain270_0(r Router, p string) { r.Get(p+"/deep0", h); chain270_1(r, p) }
func chain270_1(r Router, p string) { r.Get(p+"/deep1", h); chain270_2(r, p) }
func chain270_2(r Router, p string) { r.Get(p+"/deep2", h); chain270_3(r, p) }
func chain270_3(r Router, p string) { deep_0(r, p) }

type m271 struct{}

func (m271) Register(r Router) { reg271(r, "/m271") }
func reg271(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain271_0(r, p)
}
func chain271_0(r Router, p string) { r.Get(p+"/deep0", h); chain271_1(r, p) }
func chain271_1(r Router, p string) { r.Get(p+"/deep1", h); chain271_2(r, p) }
func chain271_2(r Router, p string) { r.Get(p+"/deep2", h); chain271_3(r, p) }
func chain271_3(r Router, p string) { deep_0(r, p) }

type m272 struct{}

func (m272) Register(r Router) { reg272(r, "/m272") }
func reg272(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain272_0(r, p)
}
func chain272_0(r Router, p string) { r.Get(p+"/deep0", h); chain272_1(r, p) }
func chain272_1(r Router, p string) { r.Get(p+"/deep1", h); chain272_2(r, p) }
func chain272_2(r Router, p string) { r.Get(p+"/deep2", h); chain272_3(r, p) }
func chain272_3(r Router, p string) { deep_0(r, p) }

type m273 struct{}

func (m273) Register(r Router) { reg273(r, "/m273") }
func reg273(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain273_0(r, p)
}
func chain273_0(r Router, p string) { r.Get(p+"/deep0", h); chain273_1(r, p) }
func chain273_1(r Router, p string) { r.Get(p+"/deep1", h); chain273_2(r, p) }
func chain273_2(r Router, p string) { r.Get(p+"/deep2", h); chain273_3(r, p) }
func chain273_3(r Router, p string) { deep_0(r, p) }

type m274 struct{}

func (m274) Register(r Router) { reg274(r, "/m274") }
func reg274(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain274_0(r, p)
}
func chain274_0(r Router, p string) { r.Get(p+"/deep0", h); chain274_1(r, p) }
func chain274_1(r Router, p string) { r.Get(p+"/deep1", h); chain274_2(r, p) }
func chain274_2(r Router, p string) { r.Get(p+"/deep2", h); chain274_3(r, p) }
func chain274_3(r Router, p string) { deep_0(r, p) }

type m275 struct{}

func (m275) Register(r Router) { reg275(r, "/m275") }
func reg275(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain275_0(r, p)
}
func chain275_0(r Router, p string) { r.Get(p+"/deep0", h); chain275_1(r, p) }
func chain275_1(r Router, p string) { r.Get(p+"/deep1", h); chain275_2(r, p) }
func chain275_2(r Router, p string) { r.Get(p+"/deep2", h); chain275_3(r, p) }
func chain275_3(r Router, p string) { deep_0(r, p) }

type m276 struct{}

func (m276) Register(r Router) { reg276(r, "/m276") }
func reg276(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain276_0(r, p)
}
func chain276_0(r Router, p string) { r.Get(p+"/deep0", h); chain276_1(r, p) }
func chain276_1(r Router, p string) { r.Get(p+"/deep1", h); chain276_2(r, p) }
func chain276_2(r Router, p string) { r.Get(p+"/deep2", h); chain276_3(r, p) }
func chain276_3(r Router, p string) { deep_0(r, p) }

type m277 struct{}

func (m277) Register(r Router) { reg277(r, "/m277") }
func reg277(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain277_0(r, p)
}
func chain277_0(r Router, p string) { r.Get(p+"/deep0", h); chain277_1(r, p) }
func chain277_1(r Router, p string) { r.Get(p+"/deep1", h); chain277_2(r, p) }
func chain277_2(r Router, p string) { r.Get(p+"/deep2", h); chain277_3(r, p) }
func chain277_3(r Router, p string) { deep_0(r, p) }

type m278 struct{}

func (m278) Register(r Router) { reg278(r, "/m278") }
func reg278(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain278_0(r, p)
}
func chain278_0(r Router, p string) { r.Get(p+"/deep0", h); chain278_1(r, p) }
func chain278_1(r Router, p string) { r.Get(p+"/deep1", h); chain278_2(r, p) }
func chain278_2(r Router, p string) { r.Get(p+"/deep2", h); chain278_3(r, p) }
func chain278_3(r Router, p string) { deep_0(r, p) }

type m279 struct{}

func (m279) Register(r Router) { reg279(r, "/m279") }
func reg279(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain279_0(r, p)
}
func chain279_0(r Router, p string) { r.Get(p+"/deep0", h); chain279_1(r, p) }
func chain279_1(r Router, p string) { r.Get(p+"/deep1", h); chain279_2(r, p) }
func chain279_2(r Router, p string) { r.Get(p+"/deep2", h); chain279_3(r, p) }
func chain279_3(r Router, p string) { deep_0(r, p) }

type m280 struct{}

func (m280) Register(r Router) { reg280(r, "/m280") }
func reg280(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain280_0(r, p)
}
func chain280_0(r Router, p string) { r.Get(p+"/deep0", h); chain280_1(r, p) }
func chain280_1(r Router, p string) { r.Get(p+"/deep1", h); chain280_2(r, p) }
func chain280_2(r Router, p string) { r.Get(p+"/deep2", h); chain280_3(r, p) }
func chain280_3(r Router, p string) { deep_0(r, p) }

type m281 struct{}

func (m281) Register(r Router) { reg281(r, "/m281") }
func reg281(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain281_0(r, p)
}
func chain281_0(r Router, p string) { r.Get(p+"/deep0", h); chain281_1(r, p) }
func chain281_1(r Router, p string) { r.Get(p+"/deep1", h); chain281_2(r, p) }
func chain281_2(r Router, p string) { r.Get(p+"/deep2", h); chain281_3(r, p) }
func chain281_3(r Router, p string) { deep_0(r, p) }

type m282 struct{}

func (m282) Register(r Router) { reg282(r, "/m282") }
func reg282(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain282_0(r, p)
}
func chain282_0(r Router, p string) { r.Get(p+"/deep0", h); chain282_1(r, p) }
func chain282_1(r Router, p string) { r.Get(p+"/deep1", h); chain282_2(r, p) }
func chain282_2(r Router, p string) { r.Get(p+"/deep2", h); chain282_3(r, p) }
func chain282_3(r Router, p string) { deep_0(r, p) }

type m283 struct{}

func (m283) Register(r Router) { reg283(r, "/m283") }
func reg283(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain283_0(r, p)
}
func chain283_0(r Router, p string) { r.Get(p+"/deep0", h); chain283_1(r, p) }
func chain283_1(r Router, p string) { r.Get(p+"/deep1", h); chain283_2(r, p) }
func chain283_2(r Router, p string) { r.Get(p+"/deep2", h); chain283_3(r, p) }
func chain283_3(r Router, p string) { deep_0(r, p) }

type m284 struct{}

func (m284) Register(r Router) { reg284(r, "/m284") }
func reg284(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain284_0(r, p)
}
func chain284_0(r Router, p string) { r.Get(p+"/deep0", h); chain284_1(r, p) }
func chain284_1(r Router, p string) { r.Get(p+"/deep1", h); chain284_2(r, p) }
func chain284_2(r Router, p string) { r.Get(p+"/deep2", h); chain284_3(r, p) }
func chain284_3(r Router, p string) { deep_0(r, p) }

type m285 struct{}

func (m285) Register(r Router) { reg285(r, "/m285") }
func reg285(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain285_0(r, p)
}
func chain285_0(r Router, p string) { r.Get(p+"/deep0", h); chain285_1(r, p) }
func chain285_1(r Router, p string) { r.Get(p+"/deep1", h); chain285_2(r, p) }
func chain285_2(r Router, p string) { r.Get(p+"/deep2", h); chain285_3(r, p) }
func chain285_3(r Router, p string) { deep_0(r, p) }

type m286 struct{}

func (m286) Register(r Router) { reg286(r, "/m286") }
func reg286(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain286_0(r, p)
}
func chain286_0(r Router, p string) { r.Get(p+"/deep0", h); chain286_1(r, p) }
func chain286_1(r Router, p string) { r.Get(p+"/deep1", h); chain286_2(r, p) }
func chain286_2(r Router, p string) { r.Get(p+"/deep2", h); chain286_3(r, p) }
func chain286_3(r Router, p string) { deep_0(r, p) }

type m287 struct{}

func (m287) Register(r Router) { reg287(r, "/m287") }
func reg287(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain287_0(r, p)
}
func chain287_0(r Router, p string) { r.Get(p+"/deep0", h); chain287_1(r, p) }
func chain287_1(r Router, p string) { r.Get(p+"/deep1", h); chain287_2(r, p) }
func chain287_2(r Router, p string) { r.Get(p+"/deep2", h); chain287_3(r, p) }
func chain287_3(r Router, p string) { deep_0(r, p) }

type m288 struct{}

func (m288) Register(r Router) { reg288(r, "/m288") }
func reg288(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain288_0(r, p)
}
func chain288_0(r Router, p string) { r.Get(p+"/deep0", h); chain288_1(r, p) }
func chain288_1(r Router, p string) { r.Get(p+"/deep1", h); chain288_2(r, p) }
func chain288_2(r Router, p string) { r.Get(p+"/deep2", h); chain288_3(r, p) }
func chain288_3(r Router, p string) { deep_0(r, p) }

type m289 struct{}

func (m289) Register(r Router) { reg289(r, "/m289") }
func reg289(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain289_0(r, p)
}
func chain289_0(r Router, p string) { r.Get(p+"/deep0", h); chain289_1(r, p) }
func chain289_1(r Router, p string) { r.Get(p+"/deep1", h); chain289_2(r, p) }
func chain289_2(r Router, p string) { r.Get(p+"/deep2", h); chain289_3(r, p) }
func chain289_3(r Router, p string) { deep_0(r, p) }

type m290 struct{}

func (m290) Register(r Router) { reg290(r, "/m290") }
func reg290(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain290_0(r, p)
}
func chain290_0(r Router, p string) { r.Get(p+"/deep0", h); chain290_1(r, p) }
func chain290_1(r Router, p string) { r.Get(p+"/deep1", h); chain290_2(r, p) }
func chain290_2(r Router, p string) { r.Get(p+"/deep2", h); chain290_3(r, p) }
func chain290_3(r Router, p string) { deep_0(r, p) }

type m291 struct{}

func (m291) Register(r Router) { reg291(r, "/m291") }
func reg291(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain291_0(r, p)
}
func chain291_0(r Router, p string) { r.Get(p+"/deep0", h); chain291_1(r, p) }
func chain291_1(r Router, p string) { r.Get(p+"/deep1", h); chain291_2(r, p) }
func chain291_2(r Router, p string) { r.Get(p+"/deep2", h); chain291_3(r, p) }
func chain291_3(r Router, p string) { deep_0(r, p) }

type m292 struct{}

func (m292) Register(r Router) { reg292(r, "/m292") }
func reg292(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain292_0(r, p)
}
func chain292_0(r Router, p string) { r.Get(p+"/deep0", h); chain292_1(r, p) }
func chain292_1(r Router, p string) { r.Get(p+"/deep1", h); chain292_2(r, p) }
func chain292_2(r Router, p string) { r.Get(p+"/deep2", h); chain292_3(r, p) }
func chain292_3(r Router, p string) { deep_0(r, p) }

type m293 struct{}

func (m293) Register(r Router) { reg293(r, "/m293") }
func reg293(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain293_0(r, p)
}
func chain293_0(r Router, p string) { r.Get(p+"/deep0", h); chain293_1(r, p) }
func chain293_1(r Router, p string) { r.Get(p+"/deep1", h); chain293_2(r, p) }
func chain293_2(r Router, p string) { r.Get(p+"/deep2", h); chain293_3(r, p) }
func chain293_3(r Router, p string) { deep_0(r, p) }

type m294 struct{}

func (m294) Register(r Router) { reg294(r, "/m294") }
func reg294(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain294_0(r, p)
}
func chain294_0(r Router, p string) { r.Get(p+"/deep0", h); chain294_1(r, p) }
func chain294_1(r Router, p string) { r.Get(p+"/deep1", h); chain294_2(r, p) }
func chain294_2(r Router, p string) { r.Get(p+"/deep2", h); chain294_3(r, p) }
func chain294_3(r Router, p string) { deep_0(r, p) }

type m295 struct{}

func (m295) Register(r Router) { reg295(r, "/m295") }
func reg295(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain295_0(r, p)
}
func chain295_0(r Router, p string) { r.Get(p+"/deep0", h); chain295_1(r, p) }
func chain295_1(r Router, p string) { r.Get(p+"/deep1", h); chain295_2(r, p) }
func chain295_2(r Router, p string) { r.Get(p+"/deep2", h); chain295_3(r, p) }
func chain295_3(r Router, p string) { deep_0(r, p) }

type m296 struct{}

func (m296) Register(r Router) { reg296(r, "/m296") }
func reg296(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain296_0(r, p)
}
func chain296_0(r Router, p string) { r.Get(p+"/deep0", h); chain296_1(r, p) }
func chain296_1(r Router, p string) { r.Get(p+"/deep1", h); chain296_2(r, p) }
func chain296_2(r Router, p string) { r.Get(p+"/deep2", h); chain296_3(r, p) }
func chain296_3(r Router, p string) { deep_0(r, p) }

type m297 struct{}

func (m297) Register(r Router) { reg297(r, "/m297") }
func reg297(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain297_0(r, p)
}
func chain297_0(r Router, p string) { r.Get(p+"/deep0", h); chain297_1(r, p) }
func chain297_1(r Router, p string) { r.Get(p+"/deep1", h); chain297_2(r, p) }
func chain297_2(r Router, p string) { r.Get(p+"/deep2", h); chain297_3(r, p) }
func chain297_3(r Router, p string) { deep_0(r, p) }

type m298 struct{}

func (m298) Register(r Router) { reg298(r, "/m298") }
func reg298(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain298_0(r, p)
}
func chain298_0(r Router, p string) { r.Get(p+"/deep0", h); chain298_1(r, p) }
func chain298_1(r Router, p string) { r.Get(p+"/deep1", h); chain298_2(r, p) }
func chain298_2(r Router, p string) { r.Get(p+"/deep2", h); chain298_3(r, p) }
func chain298_3(r Router, p string) { deep_0(r, p) }

type m299 struct{}

func (m299) Register(r Router) { reg299(r, "/m299") }
func reg299(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain299_0(r, p)
}
func chain299_0(r Router, p string) { r.Get(p+"/deep0", h); chain299_1(r, p) }
func chain299_1(r Router, p string) { r.Get(p+"/deep1", h); chain299_2(r, p) }
func chain299_2(r Router, p string) { r.Get(p+"/deep2", h); chain299_3(r, p) }
func chain299_3(r Router, p string) { deep_0(r, p) }

type m300 struct{}

func (m300) Register(r Router) { reg300(r, "/m300") }
func reg300(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain300_0(r, p)
}
func chain300_0(r Router, p string) { r.Get(p+"/deep0", h); chain300_1(r, p) }
func chain300_1(r Router, p string) { r.Get(p+"/deep1", h); chain300_2(r, p) }
func chain300_2(r Router, p string) { r.Get(p+"/deep2", h); chain300_3(r, p) }
func chain300_3(r Router, p string) { deep_0(r, p) }

type m301 struct{}

func (m301) Register(r Router) { reg301(r, "/m301") }
func reg301(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain301_0(r, p)
}
func chain301_0(r Router, p string) { r.Get(p+"/deep0", h); chain301_1(r, p) }
func chain301_1(r Router, p string) { r.Get(p+"/deep1", h); chain301_2(r, p) }
func chain301_2(r Router, p string) { r.Get(p+"/deep2", h); chain301_3(r, p) }
func chain301_3(r Router, p string) { deep_0(r, p) }

type m302 struct{}

func (m302) Register(r Router) { reg302(r, "/m302") }
func reg302(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain302_0(r, p)
}
func chain302_0(r Router, p string) { r.Get(p+"/deep0", h); chain302_1(r, p) }
func chain302_1(r Router, p string) { r.Get(p+"/deep1", h); chain302_2(r, p) }
func chain302_2(r Router, p string) { r.Get(p+"/deep2", h); chain302_3(r, p) }
func chain302_3(r Router, p string) { deep_0(r, p) }

type m303 struct{}

func (m303) Register(r Router) { reg303(r, "/m303") }
func reg303(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain303_0(r, p)
}
func chain303_0(r Router, p string) { r.Get(p+"/deep0", h); chain303_1(r, p) }
func chain303_1(r Router, p string) { r.Get(p+"/deep1", h); chain303_2(r, p) }
func chain303_2(r Router, p string) { r.Get(p+"/deep2", h); chain303_3(r, p) }
func chain303_3(r Router, p string) { deep_0(r, p) }

type m304 struct{}

func (m304) Register(r Router) { reg304(r, "/m304") }
func reg304(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain304_0(r, p)
}
func chain304_0(r Router, p string) { r.Get(p+"/deep0", h); chain304_1(r, p) }
func chain304_1(r Router, p string) { r.Get(p+"/deep1", h); chain304_2(r, p) }
func chain304_2(r Router, p string) { r.Get(p+"/deep2", h); chain304_3(r, p) }
func chain304_3(r Router, p string) { deep_0(r, p) }

type m305 struct{}

func (m305) Register(r Router) { reg305(r, "/m305") }
func reg305(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain305_0(r, p)
}
func chain305_0(r Router, p string) { r.Get(p+"/deep0", h); chain305_1(r, p) }
func chain305_1(r Router, p string) { r.Get(p+"/deep1", h); chain305_2(r, p) }
func chain305_2(r Router, p string) { r.Get(p+"/deep2", h); chain305_3(r, p) }
func chain305_3(r Router, p string) { deep_0(r, p) }

type m306 struct{}

func (m306) Register(r Router) { reg306(r, "/m306") }
func reg306(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain306_0(r, p)
}
func chain306_0(r Router, p string) { r.Get(p+"/deep0", h); chain306_1(r, p) }
func chain306_1(r Router, p string) { r.Get(p+"/deep1", h); chain306_2(r, p) }
func chain306_2(r Router, p string) { r.Get(p+"/deep2", h); chain306_3(r, p) }
func chain306_3(r Router, p string) { deep_0(r, p) }

type m307 struct{}

func (m307) Register(r Router) { reg307(r, "/m307") }
func reg307(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain307_0(r, p)
}
func chain307_0(r Router, p string) { r.Get(p+"/deep0", h); chain307_1(r, p) }
func chain307_1(r Router, p string) { r.Get(p+"/deep1", h); chain307_2(r, p) }
func chain307_2(r Router, p string) { r.Get(p+"/deep2", h); chain307_3(r, p) }
func chain307_3(r Router, p string) { deep_0(r, p) }

type m308 struct{}

func (m308) Register(r Router) { reg308(r, "/m308") }
func reg308(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain308_0(r, p)
}
func chain308_0(r Router, p string) { r.Get(p+"/deep0", h); chain308_1(r, p) }
func chain308_1(r Router, p string) { r.Get(p+"/deep1", h); chain308_2(r, p) }
func chain308_2(r Router, p string) { r.Get(p+"/deep2", h); chain308_3(r, p) }
func chain308_3(r Router, p string) { deep_0(r, p) }

type m309 struct{}

func (m309) Register(r Router) { reg309(r, "/m309") }
func reg309(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain309_0(r, p)
}
func chain309_0(r Router, p string) { r.Get(p+"/deep0", h); chain309_1(r, p) }
func chain309_1(r Router, p string) { r.Get(p+"/deep1", h); chain309_2(r, p) }
func chain309_2(r Router, p string) { r.Get(p+"/deep2", h); chain309_3(r, p) }
func chain309_3(r Router, p string) { deep_0(r, p) }

type m310 struct{}

func (m310) Register(r Router) { reg310(r, "/m310") }
func reg310(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain310_0(r, p)
}
func chain310_0(r Router, p string) { r.Get(p+"/deep0", h); chain310_1(r, p) }
func chain310_1(r Router, p string) { r.Get(p+"/deep1", h); chain310_2(r, p) }
func chain310_2(r Router, p string) { r.Get(p+"/deep2", h); chain310_3(r, p) }
func chain310_3(r Router, p string) { deep_0(r, p) }

type m311 struct{}

func (m311) Register(r Router) { reg311(r, "/m311") }
func reg311(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain311_0(r, p)
}
func chain311_0(r Router, p string) { r.Get(p+"/deep0", h); chain311_1(r, p) }
func chain311_1(r Router, p string) { r.Get(p+"/deep1", h); chain311_2(r, p) }
func chain311_2(r Router, p string) { r.Get(p+"/deep2", h); chain311_3(r, p) }
func chain311_3(r Router, p string) { deep_0(r, p) }

type m312 struct{}

func (m312) Register(r Router) { reg312(r, "/m312") }
func reg312(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain312_0(r, p)
}
func chain312_0(r Router, p string) { r.Get(p+"/deep0", h); chain312_1(r, p) }
func chain312_1(r Router, p string) { r.Get(p+"/deep1", h); chain312_2(r, p) }
func chain312_2(r Router, p string) { r.Get(p+"/deep2", h); chain312_3(r, p) }
func chain312_3(r Router, p string) { deep_0(r, p) }

type m313 struct{}

func (m313) Register(r Router) { reg313(r, "/m313") }
func reg313(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain313_0(r, p)
}
func chain313_0(r Router, p string) { r.Get(p+"/deep0", h); chain313_1(r, p) }
func chain313_1(r Router, p string) { r.Get(p+"/deep1", h); chain313_2(r, p) }
func chain313_2(r Router, p string) { r.Get(p+"/deep2", h); chain313_3(r, p) }
func chain313_3(r Router, p string) { deep_0(r, p) }

type m314 struct{}

func (m314) Register(r Router) { reg314(r, "/m314") }
func reg314(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain314_0(r, p)
}
func chain314_0(r Router, p string) { r.Get(p+"/deep0", h); chain314_1(r, p) }
func chain314_1(r Router, p string) { r.Get(p+"/deep1", h); chain314_2(r, p) }
func chain314_2(r Router, p string) { r.Get(p+"/deep2", h); chain314_3(r, p) }
func chain314_3(r Router, p string) { deep_0(r, p) }

type m315 struct{}

func (m315) Register(r Router) { reg315(r, "/m315") }
func reg315(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain315_0(r, p)
}
func chain315_0(r Router, p string) { r.Get(p+"/deep0", h); chain315_1(r, p) }
func chain315_1(r Router, p string) { r.Get(p+"/deep1", h); chain315_2(r, p) }
func chain315_2(r Router, p string) { r.Get(p+"/deep2", h); chain315_3(r, p) }
func chain315_3(r Router, p string) { deep_0(r, p) }

type m316 struct{}

func (m316) Register(r Router) { reg316(r, "/m316") }
func reg316(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain316_0(r, p)
}
func chain316_0(r Router, p string) { r.Get(p+"/deep0", h); chain316_1(r, p) }
func chain316_1(r Router, p string) { r.Get(p+"/deep1", h); chain316_2(r, p) }
func chain316_2(r Router, p string) { r.Get(p+"/deep2", h); chain316_3(r, p) }
func chain316_3(r Router, p string) { deep_0(r, p) }

type m317 struct{}

func (m317) Register(r Router) { reg317(r, "/m317") }
func reg317(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain317_0(r, p)
}
func chain317_0(r Router, p string) { r.Get(p+"/deep0", h); chain317_1(r, p) }
func chain317_1(r Router, p string) { r.Get(p+"/deep1", h); chain317_2(r, p) }
func chain317_2(r Router, p string) { r.Get(p+"/deep2", h); chain317_3(r, p) }
func chain317_3(r Router, p string) { deep_0(r, p) }

type m318 struct{}

func (m318) Register(r Router) { reg318(r, "/m318") }
func reg318(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain318_0(r, p)
}
func chain318_0(r Router, p string) { r.Get(p+"/deep0", h); chain318_1(r, p) }
func chain318_1(r Router, p string) { r.Get(p+"/deep1", h); chain318_2(r, p) }
func chain318_2(r Router, p string) { r.Get(p+"/deep2", h); chain318_3(r, p) }
func chain318_3(r Router, p string) { deep_0(r, p) }

type m319 struct{}

func (m319) Register(r Router) { reg319(r, "/m319") }
func reg319(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain319_0(r, p)
}
func chain319_0(r Router, p string) { r.Get(p+"/deep0", h); chain319_1(r, p) }
func chain319_1(r Router, p string) { r.Get(p+"/deep1", h); chain319_2(r, p) }
func chain319_2(r Router, p string) { r.Get(p+"/deep2", h); chain319_3(r, p) }
func chain319_3(r Router, p string) { deep_0(r, p) }

type m320 struct{}

func (m320) Register(r Router) { reg320(r, "/m320") }
func reg320(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain320_0(r, p)
}
func chain320_0(r Router, p string) { r.Get(p+"/deep0", h); chain320_1(r, p) }
func chain320_1(r Router, p string) { r.Get(p+"/deep1", h); chain320_2(r, p) }
func chain320_2(r Router, p string) { r.Get(p+"/deep2", h); chain320_3(r, p) }
func chain320_3(r Router, p string) { deep_0(r, p) }

type m321 struct{}

func (m321) Register(r Router) { reg321(r, "/m321") }
func reg321(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain321_0(r, p)
}
func chain321_0(r Router, p string) { r.Get(p+"/deep0", h); chain321_1(r, p) }
func chain321_1(r Router, p string) { r.Get(p+"/deep1", h); chain321_2(r, p) }
func chain321_2(r Router, p string) { r.Get(p+"/deep2", h); chain321_3(r, p) }
func chain321_3(r Router, p string) { deep_0(r, p) }

type m322 struct{}

func (m322) Register(r Router) { reg322(r, "/m322") }
func reg322(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain322_0(r, p)
}
func chain322_0(r Router, p string) { r.Get(p+"/deep0", h); chain322_1(r, p) }
func chain322_1(r Router, p string) { r.Get(p+"/deep1", h); chain322_2(r, p) }
func chain322_2(r Router, p string) { r.Get(p+"/deep2", h); chain322_3(r, p) }
func chain322_3(r Router, p string) { deep_0(r, p) }

type m323 struct{}

func (m323) Register(r Router) { reg323(r, "/m323") }
func reg323(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain323_0(r, p)
}
func chain323_0(r Router, p string) { r.Get(p+"/deep0", h); chain323_1(r, p) }
func chain323_1(r Router, p string) { r.Get(p+"/deep1", h); chain323_2(r, p) }
func chain323_2(r Router, p string) { r.Get(p+"/deep2", h); chain323_3(r, p) }
func chain323_3(r Router, p string) { deep_0(r, p) }

type m324 struct{}

func (m324) Register(r Router) { reg324(r, "/m324") }
func reg324(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain324_0(r, p)
}
func chain324_0(r Router, p string) { r.Get(p+"/deep0", h); chain324_1(r, p) }
func chain324_1(r Router, p string) { r.Get(p+"/deep1", h); chain324_2(r, p) }
func chain324_2(r Router, p string) { r.Get(p+"/deep2", h); chain324_3(r, p) }
func chain324_3(r Router, p string) { deep_0(r, p) }

type m325 struct{}

func (m325) Register(r Router) { reg325(r, "/m325") }
func reg325(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain325_0(r, p)
}
func chain325_0(r Router, p string) { r.Get(p+"/deep0", h); chain325_1(r, p) }
func chain325_1(r Router, p string) { r.Get(p+"/deep1", h); chain325_2(r, p) }
func chain325_2(r Router, p string) { r.Get(p+"/deep2", h); chain325_3(r, p) }
func chain325_3(r Router, p string) { deep_0(r, p) }

type m326 struct{}

func (m326) Register(r Router) { reg326(r, "/m326") }
func reg326(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain326_0(r, p)
}
func chain326_0(r Router, p string) { r.Get(p+"/deep0", h); chain326_1(r, p) }
func chain326_1(r Router, p string) { r.Get(p+"/deep1", h); chain326_2(r, p) }
func chain326_2(r Router, p string) { r.Get(p+"/deep2", h); chain326_3(r, p) }
func chain326_3(r Router, p string) { deep_0(r, p) }

type m327 struct{}

func (m327) Register(r Router) { reg327(r, "/m327") }
func reg327(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain327_0(r, p)
}
func chain327_0(r Router, p string) { r.Get(p+"/deep0", h); chain327_1(r, p) }
func chain327_1(r Router, p string) { r.Get(p+"/deep1", h); chain327_2(r, p) }
func chain327_2(r Router, p string) { r.Get(p+"/deep2", h); chain327_3(r, p) }
func chain327_3(r Router, p string) { deep_0(r, p) }

type m328 struct{}

func (m328) Register(r Router) { reg328(r, "/m328") }
func reg328(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain328_0(r, p)
}
func chain328_0(r Router, p string) { r.Get(p+"/deep0", h); chain328_1(r, p) }
func chain328_1(r Router, p string) { r.Get(p+"/deep1", h); chain328_2(r, p) }
func chain328_2(r Router, p string) { r.Get(p+"/deep2", h); chain328_3(r, p) }
func chain328_3(r Router, p string) { deep_0(r, p) }

type m329 struct{}

func (m329) Register(r Router) { reg329(r, "/m329") }
func reg329(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain329_0(r, p)
}
func chain329_0(r Router, p string) { r.Get(p+"/deep0", h); chain329_1(r, p) }
func chain329_1(r Router, p string) { r.Get(p+"/deep1", h); chain329_2(r, p) }
func chain329_2(r Router, p string) { r.Get(p+"/deep2", h); chain329_3(r, p) }
func chain329_3(r Router, p string) { deep_0(r, p) }

type m330 struct{}

func (m330) Register(r Router) { reg330(r, "/m330") }
func reg330(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain330_0(r, p)
}
func chain330_0(r Router, p string) { r.Get(p+"/deep0", h); chain330_1(r, p) }
func chain330_1(r Router, p string) { r.Get(p+"/deep1", h); chain330_2(r, p) }
func chain330_2(r Router, p string) { r.Get(p+"/deep2", h); chain330_3(r, p) }
func chain330_3(r Router, p string) { deep_0(r, p) }

type m331 struct{}

func (m331) Register(r Router) { reg331(r, "/m331") }
func reg331(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain331_0(r, p)
}
func chain331_0(r Router, p string) { r.Get(p+"/deep0", h); chain331_1(r, p) }
func chain331_1(r Router, p string) { r.Get(p+"/deep1", h); chain331_2(r, p) }
func chain331_2(r Router, p string) { r.Get(p+"/deep2", h); chain331_3(r, p) }
func chain331_3(r Router, p string) { deep_0(r, p) }

type m332 struct{}

func (m332) Register(r Router) { reg332(r, "/m332") }
func reg332(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain332_0(r, p)
}
func chain332_0(r Router, p string) { r.Get(p+"/deep0", h); chain332_1(r, p) }
func chain332_1(r Router, p string) { r.Get(p+"/deep1", h); chain332_2(r, p) }
func chain332_2(r Router, p string) { r.Get(p+"/deep2", h); chain332_3(r, p) }
func chain332_3(r Router, p string) { deep_0(r, p) }

type m333 struct{}

func (m333) Register(r Router) { reg333(r, "/m333") }
func reg333(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain333_0(r, p)
}
func chain333_0(r Router, p string) { r.Get(p+"/deep0", h); chain333_1(r, p) }
func chain333_1(r Router, p string) { r.Get(p+"/deep1", h); chain333_2(r, p) }
func chain333_2(r Router, p string) { r.Get(p+"/deep2", h); chain333_3(r, p) }
func chain333_3(r Router, p string) { deep_0(r, p) }

type m334 struct{}

func (m334) Register(r Router) { reg334(r, "/m334") }
func reg334(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain334_0(r, p)
}
func chain334_0(r Router, p string) { r.Get(p+"/deep0", h); chain334_1(r, p) }
func chain334_1(r Router, p string) { r.Get(p+"/deep1", h); chain334_2(r, p) }
func chain334_2(r Router, p string) { r.Get(p+"/deep2", h); chain334_3(r, p) }
func chain334_3(r Router, p string) { deep_0(r, p) }

type m335 struct{}

func (m335) Register(r Router) { reg335(r, "/m335") }
func reg335(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain335_0(r, p)
}
func chain335_0(r Router, p string) { r.Get(p+"/deep0", h); chain335_1(r, p) }
func chain335_1(r Router, p string) { r.Get(p+"/deep1", h); chain335_2(r, p) }
func chain335_2(r Router, p string) { r.Get(p+"/deep2", h); chain335_3(r, p) }
func chain335_3(r Router, p string) { deep_0(r, p) }

type m336 struct{}

func (m336) Register(r Router) { reg336(r, "/m336") }
func reg336(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain336_0(r, p)
}
func chain336_0(r Router, p string) { r.Get(p+"/deep0", h); chain336_1(r, p) }
func chain336_1(r Router, p string) { r.Get(p+"/deep1", h); chain336_2(r, p) }
func chain336_2(r Router, p string) { r.Get(p+"/deep2", h); chain336_3(r, p) }
func chain336_3(r Router, p string) { deep_0(r, p) }

type m337 struct{}

func (m337) Register(r Router) { reg337(r, "/m337") }
func reg337(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain337_0(r, p)
}
func chain337_0(r Router, p string) { r.Get(p+"/deep0", h); chain337_1(r, p) }
func chain337_1(r Router, p string) { r.Get(p+"/deep1", h); chain337_2(r, p) }
func chain337_2(r Router, p string) { r.Get(p+"/deep2", h); chain337_3(r, p) }
func chain337_3(r Router, p string) { deep_0(r, p) }

type m338 struct{}

func (m338) Register(r Router) { reg338(r, "/m338") }
func reg338(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain338_0(r, p)
}
func chain338_0(r Router, p string) { r.Get(p+"/deep0", h); chain338_1(r, p) }
func chain338_1(r Router, p string) { r.Get(p+"/deep1", h); chain338_2(r, p) }
func chain338_2(r Router, p string) { r.Get(p+"/deep2", h); chain338_3(r, p) }
func chain338_3(r Router, p string) { deep_0(r, p) }

type m339 struct{}

func (m339) Register(r Router) { reg339(r, "/m339") }
func reg339(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain339_0(r, p)
}
func chain339_0(r Router, p string) { r.Get(p+"/deep0", h); chain339_1(r, p) }
func chain339_1(r Router, p string) { r.Get(p+"/deep1", h); chain339_2(r, p) }
func chain339_2(r Router, p string) { r.Get(p+"/deep2", h); chain339_3(r, p) }
func chain339_3(r Router, p string) { deep_0(r, p) }

type m340 struct{}

func (m340) Register(r Router) { reg340(r, "/m340") }
func reg340(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain340_0(r, p)
}
func chain340_0(r Router, p string) { r.Get(p+"/deep0", h); chain340_1(r, p) }
func chain340_1(r Router, p string) { r.Get(p+"/deep1", h); chain340_2(r, p) }
func chain340_2(r Router, p string) { r.Get(p+"/deep2", h); chain340_3(r, p) }
func chain340_3(r Router, p string) { deep_0(r, p) }

type m341 struct{}

func (m341) Register(r Router) { reg341(r, "/m341") }
func reg341(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain341_0(r, p)
}
func chain341_0(r Router, p string) { r.Get(p+"/deep0", h); chain341_1(r, p) }
func chain341_1(r Router, p string) { r.Get(p+"/deep1", h); chain341_2(r, p) }
func chain341_2(r Router, p string) { r.Get(p+"/deep2", h); chain341_3(r, p) }
func chain341_3(r Router, p string) { deep_0(r, p) }

type m342 struct{}

func (m342) Register(r Router) { reg342(r, "/m342") }
func reg342(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain342_0(r, p)
}
func chain342_0(r Router, p string) { r.Get(p+"/deep0", h); chain342_1(r, p) }
func chain342_1(r Router, p string) { r.Get(p+"/deep1", h); chain342_2(r, p) }
func chain342_2(r Router, p string) { r.Get(p+"/deep2", h); chain342_3(r, p) }
func chain342_3(r Router, p string) { deep_0(r, p) }

type m343 struct{}

func (m343) Register(r Router) { reg343(r, "/m343") }
func reg343(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain343_0(r, p)
}
func chain343_0(r Router, p string) { r.Get(p+"/deep0", h); chain343_1(r, p) }
func chain343_1(r Router, p string) { r.Get(p+"/deep1", h); chain343_2(r, p) }
func chain343_2(r Router, p string) { r.Get(p+"/deep2", h); chain343_3(r, p) }
func chain343_3(r Router, p string) { deep_0(r, p) }

type m344 struct{}

func (m344) Register(r Router) { reg344(r, "/m344") }
func reg344(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain344_0(r, p)
}
func chain344_0(r Router, p string) { r.Get(p+"/deep0", h); chain344_1(r, p) }
func chain344_1(r Router, p string) { r.Get(p+"/deep1", h); chain344_2(r, p) }
func chain344_2(r Router, p string) { r.Get(p+"/deep2", h); chain344_3(r, p) }
func chain344_3(r Router, p string) { deep_0(r, p) }

type m345 struct{}

func (m345) Register(r Router) { reg345(r, "/m345") }
func reg345(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain345_0(r, p)
}
func chain345_0(r Router, p string) { r.Get(p+"/deep0", h); chain345_1(r, p) }
func chain345_1(r Router, p string) { r.Get(p+"/deep1", h); chain345_2(r, p) }
func chain345_2(r Router, p string) { r.Get(p+"/deep2", h); chain345_3(r, p) }
func chain345_3(r Router, p string) { deep_0(r, p) }

type m346 struct{}

func (m346) Register(r Router) { reg346(r, "/m346") }
func reg346(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain346_0(r, p)
}
func chain346_0(r Router, p string) { r.Get(p+"/deep0", h); chain346_1(r, p) }
func chain346_1(r Router, p string) { r.Get(p+"/deep1", h); chain346_2(r, p) }
func chain346_2(r Router, p string) { r.Get(p+"/deep2", h); chain346_3(r, p) }
func chain346_3(r Router, p string) { deep_0(r, p) }

type m347 struct{}

func (m347) Register(r Router) { reg347(r, "/m347") }
func reg347(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain347_0(r, p)
}
func chain347_0(r Router, p string) { r.Get(p+"/deep0", h); chain347_1(r, p) }
func chain347_1(r Router, p string) { r.Get(p+"/deep1", h); chain347_2(r, p) }
func chain347_2(r Router, p string) { r.Get(p+"/deep2", h); chain347_3(r, p) }
func chain347_3(r Router, p string) { deep_0(r, p) }

type m348 struct{}

func (m348) Register(r Router) { reg348(r, "/m348") }
func reg348(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain348_0(r, p)
}
func chain348_0(r Router, p string) { r.Get(p+"/deep0", h); chain348_1(r, p) }
func chain348_1(r Router, p string) { r.Get(p+"/deep1", h); chain348_2(r, p) }
func chain348_2(r Router, p string) { r.Get(p+"/deep2", h); chain348_3(r, p) }
func chain348_3(r Router, p string) { deep_0(r, p) }

type m349 struct{}

func (m349) Register(r Router) { reg349(r, "/m349") }
func reg349(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain349_0(r, p)
}
func chain349_0(r Router, p string) { r.Get(p+"/deep0", h); chain349_1(r, p) }
func chain349_1(r Router, p string) { r.Get(p+"/deep1", h); chain349_2(r, p) }
func chain349_2(r Router, p string) { r.Get(p+"/deep2", h); chain349_3(r, p) }
func chain349_3(r Router, p string) { deep_0(r, p) }

type m350 struct{}

func (m350) Register(r Router) { reg350(r, "/m350") }
func reg350(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain350_0(r, p)
}
func chain350_0(r Router, p string) { r.Get(p+"/deep0", h); chain350_1(r, p) }
func chain350_1(r Router, p string) { r.Get(p+"/deep1", h); chain350_2(r, p) }
func chain350_2(r Router, p string) { r.Get(p+"/deep2", h); chain350_3(r, p) }
func chain350_3(r Router, p string) { deep_0(r, p) }

type m351 struct{}

func (m351) Register(r Router) { reg351(r, "/m351") }
func reg351(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain351_0(r, p)
}
func chain351_0(r Router, p string) { r.Get(p+"/deep0", h); chain351_1(r, p) }
func chain351_1(r Router, p string) { r.Get(p+"/deep1", h); chain351_2(r, p) }
func chain351_2(r Router, p string) { r.Get(p+"/deep2", h); chain351_3(r, p) }
func chain351_3(r Router, p string) { deep_0(r, p) }

type m352 struct{}

func (m352) Register(r Router) { reg352(r, "/m352") }
func reg352(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain352_0(r, p)
}
func chain352_0(r Router, p string) { r.Get(p+"/deep0", h); chain352_1(r, p) }
func chain352_1(r Router, p string) { r.Get(p+"/deep1", h); chain352_2(r, p) }
func chain352_2(r Router, p string) { r.Get(p+"/deep2", h); chain352_3(r, p) }
func chain352_3(r Router, p string) { deep_0(r, p) }

type m353 struct{}

func (m353) Register(r Router) { reg353(r, "/m353") }
func reg353(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain353_0(r, p)
}
func chain353_0(r Router, p string) { r.Get(p+"/deep0", h); chain353_1(r, p) }
func chain353_1(r Router, p string) { r.Get(p+"/deep1", h); chain353_2(r, p) }
func chain353_2(r Router, p string) { r.Get(p+"/deep2", h); chain353_3(r, p) }
func chain353_3(r Router, p string) { deep_0(r, p) }

type m354 struct{}

func (m354) Register(r Router) { reg354(r, "/m354") }
func reg354(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain354_0(r, p)
}
func chain354_0(r Router, p string) { r.Get(p+"/deep0", h); chain354_1(r, p) }
func chain354_1(r Router, p string) { r.Get(p+"/deep1", h); chain354_2(r, p) }
func chain354_2(r Router, p string) { r.Get(p+"/deep2", h); chain354_3(r, p) }
func chain354_3(r Router, p string) { deep_0(r, p) }

type m355 struct{}

func (m355) Register(r Router) { reg355(r, "/m355") }
func reg355(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain355_0(r, p)
}
func chain355_0(r Router, p string) { r.Get(p+"/deep0", h); chain355_1(r, p) }
func chain355_1(r Router, p string) { r.Get(p+"/deep1", h); chain355_2(r, p) }
func chain355_2(r Router, p string) { r.Get(p+"/deep2", h); chain355_3(r, p) }
func chain355_3(r Router, p string) { deep_0(r, p) }

type m356 struct{}

func (m356) Register(r Router) { reg356(r, "/m356") }
func reg356(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain356_0(r, p)
}
func chain356_0(r Router, p string) { r.Get(p+"/deep0", h); chain356_1(r, p) }
func chain356_1(r Router, p string) { r.Get(p+"/deep1", h); chain356_2(r, p) }
func chain356_2(r Router, p string) { r.Get(p+"/deep2", h); chain356_3(r, p) }
func chain356_3(r Router, p string) { deep_0(r, p) }

type m357 struct{}

func (m357) Register(r Router) { reg357(r, "/m357") }
func reg357(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain357_0(r, p)
}
func chain357_0(r Router, p string) { r.Get(p+"/deep0", h); chain357_1(r, p) }
func chain357_1(r Router, p string) { r.Get(p+"/deep1", h); chain357_2(r, p) }
func chain357_2(r Router, p string) { r.Get(p+"/deep2", h); chain357_3(r, p) }
func chain357_3(r Router, p string) { deep_0(r, p) }

type m358 struct{}

func (m358) Register(r Router) { reg358(r, "/m358") }
func reg358(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain358_0(r, p)
}
func chain358_0(r Router, p string) { r.Get(p+"/deep0", h); chain358_1(r, p) }
func chain358_1(r Router, p string) { r.Get(p+"/deep1", h); chain358_2(r, p) }
func chain358_2(r Router, p string) { r.Get(p+"/deep2", h); chain358_3(r, p) }
func chain358_3(r Router, p string) { deep_0(r, p) }

type m359 struct{}

func (m359) Register(r Router) { reg359(r, "/m359") }
func reg359(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain359_0(r, p)
}
func chain359_0(r Router, p string) { r.Get(p+"/deep0", h); chain359_1(r, p) }
func chain359_1(r Router, p string) { r.Get(p+"/deep1", h); chain359_2(r, p) }
func chain359_2(r Router, p string) { r.Get(p+"/deep2", h); chain359_3(r, p) }
func chain359_3(r Router, p string) { deep_0(r, p) }

type m360 struct{}

func (m360) Register(r Router) { reg360(r, "/m360") }
func reg360(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain360_0(r, p)
}
func chain360_0(r Router, p string) { r.Get(p+"/deep0", h); chain360_1(r, p) }
func chain360_1(r Router, p string) { r.Get(p+"/deep1", h); chain360_2(r, p) }
func chain360_2(r Router, p string) { r.Get(p+"/deep2", h); chain360_3(r, p) }
func chain360_3(r Router, p string) { deep_0(r, p) }

type m361 struct{}

func (m361) Register(r Router) { reg361(r, "/m361") }
func reg361(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain361_0(r, p)
}
func chain361_0(r Router, p string) { r.Get(p+"/deep0", h); chain361_1(r, p) }
func chain361_1(r Router, p string) { r.Get(p+"/deep1", h); chain361_2(r, p) }
func chain361_2(r Router, p string) { r.Get(p+"/deep2", h); chain361_3(r, p) }
func chain361_3(r Router, p string) { deep_0(r, p) }

type m362 struct{}

func (m362) Register(r Router) { reg362(r, "/m362") }
func reg362(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain362_0(r, p)
}
func chain362_0(r Router, p string) { r.Get(p+"/deep0", h); chain362_1(r, p) }
func chain362_1(r Router, p string) { r.Get(p+"/deep1", h); chain362_2(r, p) }
func chain362_2(r Router, p string) { r.Get(p+"/deep2", h); chain362_3(r, p) }
func chain362_3(r Router, p string) { deep_0(r, p) }

type m363 struct{}

func (m363) Register(r Router) { reg363(r, "/m363") }
func reg363(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain363_0(r, p)
}
func chain363_0(r Router, p string) { r.Get(p+"/deep0", h); chain363_1(r, p) }
func chain363_1(r Router, p string) { r.Get(p+"/deep1", h); chain363_2(r, p) }
func chain363_2(r Router, p string) { r.Get(p+"/deep2", h); chain363_3(r, p) }
func chain363_3(r Router, p string) { deep_0(r, p) }

type m364 struct{}

func (m364) Register(r Router) { reg364(r, "/m364") }
func reg364(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain364_0(r, p)
}
func chain364_0(r Router, p string) { r.Get(p+"/deep0", h); chain364_1(r, p) }
func chain364_1(r Router, p string) { r.Get(p+"/deep1", h); chain364_2(r, p) }
func chain364_2(r Router, p string) { r.Get(p+"/deep2", h); chain364_3(r, p) }
func chain364_3(r Router, p string) { deep_0(r, p) }

type m365 struct{}

func (m365) Register(r Router) { reg365(r, "/m365") }
func reg365(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain365_0(r, p)
}
func chain365_0(r Router, p string) { r.Get(p+"/deep0", h); chain365_1(r, p) }
func chain365_1(r Router, p string) { r.Get(p+"/deep1", h); chain365_2(r, p) }
func chain365_2(r Router, p string) { r.Get(p+"/deep2", h); chain365_3(r, p) }
func chain365_3(r Router, p string) { deep_0(r, p) }

type m366 struct{}

func (m366) Register(r Router) { reg366(r, "/m366") }
func reg366(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain366_0(r, p)
}
func chain366_0(r Router, p string) { r.Get(p+"/deep0", h); chain366_1(r, p) }
func chain366_1(r Router, p string) { r.Get(p+"/deep1", h); chain366_2(r, p) }
func chain366_2(r Router, p string) { r.Get(p+"/deep2", h); chain366_3(r, p) }
func chain366_3(r Router, p string) { deep_0(r, p) }

type m367 struct{}

func (m367) Register(r Router) { reg367(r, "/m367") }
func reg367(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain367_0(r, p)
}
func chain367_0(r Router, p string) { r.Get(p+"/deep0", h); chain367_1(r, p) }
func chain367_1(r Router, p string) { r.Get(p+"/deep1", h); chain367_2(r, p) }
func chain367_2(r Router, p string) { r.Get(p+"/deep2", h); chain367_3(r, p) }
func chain367_3(r Router, p string) { deep_0(r, p) }

type m368 struct{}

func (m368) Register(r Router) { reg368(r, "/m368") }
func reg368(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain368_0(r, p)
}
func chain368_0(r Router, p string) { r.Get(p+"/deep0", h); chain368_1(r, p) }
func chain368_1(r Router, p string) { r.Get(p+"/deep1", h); chain368_2(r, p) }
func chain368_2(r Router, p string) { r.Get(p+"/deep2", h); chain368_3(r, p) }
func chain368_3(r Router, p string) { deep_0(r, p) }

type m369 struct{}

func (m369) Register(r Router) { reg369(r, "/m369") }
func reg369(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain369_0(r, p)
}
func chain369_0(r Router, p string) { r.Get(p+"/deep0", h); chain369_1(r, p) }
func chain369_1(r Router, p string) { r.Get(p+"/deep1", h); chain369_2(r, p) }
func chain369_2(r Router, p string) { r.Get(p+"/deep2", h); chain369_3(r, p) }
func chain369_3(r Router, p string) { deep_0(r, p) }

type m370 struct{}

func (m370) Register(r Router) { reg370(r, "/m370") }
func reg370(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain370_0(r, p)
}
func chain370_0(r Router, p string) { r.Get(p+"/deep0", h); chain370_1(r, p) }
func chain370_1(r Router, p string) { r.Get(p+"/deep1", h); chain370_2(r, p) }
func chain370_2(r Router, p string) { r.Get(p+"/deep2", h); chain370_3(r, p) }
func chain370_3(r Router, p string) { deep_0(r, p) }

type m371 struct{}

func (m371) Register(r Router) { reg371(r, "/m371") }
func reg371(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain371_0(r, p)
}
func chain371_0(r Router, p string) { r.Get(p+"/deep0", h); chain371_1(r, p) }
func chain371_1(r Router, p string) { r.Get(p+"/deep1", h); chain371_2(r, p) }
func chain371_2(r Router, p string) { r.Get(p+"/deep2", h); chain371_3(r, p) }
func chain371_3(r Router, p string) { deep_0(r, p) }

type m372 struct{}

func (m372) Register(r Router) { reg372(r, "/m372") }
func reg372(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain372_0(r, p)
}
func chain372_0(r Router, p string) { r.Get(p+"/deep0", h); chain372_1(r, p) }
func chain372_1(r Router, p string) { r.Get(p+"/deep1", h); chain372_2(r, p) }
func chain372_2(r Router, p string) { r.Get(p+"/deep2", h); chain372_3(r, p) }
func chain372_3(r Router, p string) { deep_0(r, p) }

type m373 struct{}

func (m373) Register(r Router) { reg373(r, "/m373") }
func reg373(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain373_0(r, p)
}
func chain373_0(r Router, p string) { r.Get(p+"/deep0", h); chain373_1(r, p) }
func chain373_1(r Router, p string) { r.Get(p+"/deep1", h); chain373_2(r, p) }
func chain373_2(r Router, p string) { r.Get(p+"/deep2", h); chain373_3(r, p) }
func chain373_3(r Router, p string) { deep_0(r, p) }

type m374 struct{}

func (m374) Register(r Router) { reg374(r, "/m374") }
func reg374(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain374_0(r, p)
}
func chain374_0(r Router, p string) { r.Get(p+"/deep0", h); chain374_1(r, p) }
func chain374_1(r Router, p string) { r.Get(p+"/deep1", h); chain374_2(r, p) }
func chain374_2(r Router, p string) { r.Get(p+"/deep2", h); chain374_3(r, p) }
func chain374_3(r Router, p string) { deep_0(r, p) }

type m375 struct{}

func (m375) Register(r Router) { reg375(r, "/m375") }
func reg375(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain375_0(r, p)
}
func chain375_0(r Router, p string) { r.Get(p+"/deep0", h); chain375_1(r, p) }
func chain375_1(r Router, p string) { r.Get(p+"/deep1", h); chain375_2(r, p) }
func chain375_2(r Router, p string) { r.Get(p+"/deep2", h); chain375_3(r, p) }
func chain375_3(r Router, p string) { deep_0(r, p) }

type m376 struct{}

func (m376) Register(r Router) { reg376(r, "/m376") }
func reg376(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain376_0(r, p)
}
func chain376_0(r Router, p string) { r.Get(p+"/deep0", h); chain376_1(r, p) }
func chain376_1(r Router, p string) { r.Get(p+"/deep1", h); chain376_2(r, p) }
func chain376_2(r Router, p string) { r.Get(p+"/deep2", h); chain376_3(r, p) }
func chain376_3(r Router, p string) { deep_0(r, p) }

type m377 struct{}

func (m377) Register(r Router) { reg377(r, "/m377") }
func reg377(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain377_0(r, p)
}
func chain377_0(r Router, p string) { r.Get(p+"/deep0", h); chain377_1(r, p) }
func chain377_1(r Router, p string) { r.Get(p+"/deep1", h); chain377_2(r, p) }
func chain377_2(r Router, p string) { r.Get(p+"/deep2", h); chain377_3(r, p) }
func chain377_3(r Router, p string) { deep_0(r, p) }

type m378 struct{}

func (m378) Register(r Router) { reg378(r, "/m378") }
func reg378(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain378_0(r, p)
}
func chain378_0(r Router, p string) { r.Get(p+"/deep0", h); chain378_1(r, p) }
func chain378_1(r Router, p string) { r.Get(p+"/deep1", h); chain378_2(r, p) }
func chain378_2(r Router, p string) { r.Get(p+"/deep2", h); chain378_3(r, p) }
func chain378_3(r Router, p string) { deep_0(r, p) }

type m379 struct{}

func (m379) Register(r Router) { reg379(r, "/m379") }
func reg379(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain379_0(r, p)
}
func chain379_0(r Router, p string) { r.Get(p+"/deep0", h); chain379_1(r, p) }
func chain379_1(r Router, p string) { r.Get(p+"/deep1", h); chain379_2(r, p) }
func chain379_2(r Router, p string) { r.Get(p+"/deep2", h); chain379_3(r, p) }
func chain379_3(r Router, p string) { deep_0(r, p) }

type m380 struct{}

func (m380) Register(r Router) { reg380(r, "/m380") }
func reg380(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain380_0(r, p)
}
func chain380_0(r Router, p string) { r.Get(p+"/deep0", h); chain380_1(r, p) }
func chain380_1(r Router, p string) { r.Get(p+"/deep1", h); chain380_2(r, p) }
func chain380_2(r Router, p string) { r.Get(p+"/deep2", h); chain380_3(r, p) }
func chain380_3(r Router, p string) { deep_0(r, p) }

type m381 struct{}

func (m381) Register(r Router) { reg381(r, "/m381") }
func reg381(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain381_0(r, p)
}
func chain381_0(r Router, p string) { r.Get(p+"/deep0", h); chain381_1(r, p) }
func chain381_1(r Router, p string) { r.Get(p+"/deep1", h); chain381_2(r, p) }
func chain381_2(r Router, p string) { r.Get(p+"/deep2", h); chain381_3(r, p) }
func chain381_3(r Router, p string) { deep_0(r, p) }

type m382 struct{}

func (m382) Register(r Router) { reg382(r, "/m382") }
func reg382(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain382_0(r, p)
}
func chain382_0(r Router, p string) { r.Get(p+"/deep0", h); chain382_1(r, p) }
func chain382_1(r Router, p string) { r.Get(p+"/deep1", h); chain382_2(r, p) }
func chain382_2(r Router, p string) { r.Get(p+"/deep2", h); chain382_3(r, p) }
func chain382_3(r Router, p string) { deep_0(r, p) }

type m383 struct{}

func (m383) Register(r Router) { reg383(r, "/m383") }
func reg383(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain383_0(r, p)
}
func chain383_0(r Router, p string) { r.Get(p+"/deep0", h); chain383_1(r, p) }
func chain383_1(r Router, p string) { r.Get(p+"/deep1", h); chain383_2(r, p) }
func chain383_2(r Router, p string) { r.Get(p+"/deep2", h); chain383_3(r, p) }
func chain383_3(r Router, p string) { deep_0(r, p) }

type m384 struct{}

func (m384) Register(r Router) { reg384(r, "/m384") }
func reg384(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain384_0(r, p)
}
func chain384_0(r Router, p string) { r.Get(p+"/deep0", h); chain384_1(r, p) }
func chain384_1(r Router, p string) { r.Get(p+"/deep1", h); chain384_2(r, p) }
func chain384_2(r Router, p string) { r.Get(p+"/deep2", h); chain384_3(r, p) }
func chain384_3(r Router, p string) { deep_0(r, p) }

type m385 struct{}

func (m385) Register(r Router) { reg385(r, "/m385") }
func reg385(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain385_0(r, p)
}
func chain385_0(r Router, p string) { r.Get(p+"/deep0", h); chain385_1(r, p) }
func chain385_1(r Router, p string) { r.Get(p+"/deep1", h); chain385_2(r, p) }
func chain385_2(r Router, p string) { r.Get(p+"/deep2", h); chain385_3(r, p) }
func chain385_3(r Router, p string) { deep_0(r, p) }

type m386 struct{}

func (m386) Register(r Router) { reg386(r, "/m386") }
func reg386(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain386_0(r, p)
}
func chain386_0(r Router, p string) { r.Get(p+"/deep0", h); chain386_1(r, p) }
func chain386_1(r Router, p string) { r.Get(p+"/deep1", h); chain386_2(r, p) }
func chain386_2(r Router, p string) { r.Get(p+"/deep2", h); chain386_3(r, p) }
func chain386_3(r Router, p string) { deep_0(r, p) }

type m387 struct{}

func (m387) Register(r Router) { reg387(r, "/m387") }
func reg387(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain387_0(r, p)
}
func chain387_0(r Router, p string) { r.Get(p+"/deep0", h); chain387_1(r, p) }
func chain387_1(r Router, p string) { r.Get(p+"/deep1", h); chain387_2(r, p) }
func chain387_2(r Router, p string) { r.Get(p+"/deep2", h); chain387_3(r, p) }
func chain387_3(r Router, p string) { deep_0(r, p) }

type m388 struct{}

func (m388) Register(r Router) { reg388(r, "/m388") }
func reg388(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain388_0(r, p)
}
func chain388_0(r Router, p string) { r.Get(p+"/deep0", h); chain388_1(r, p) }
func chain388_1(r Router, p string) { r.Get(p+"/deep1", h); chain388_2(r, p) }
func chain388_2(r Router, p string) { r.Get(p+"/deep2", h); chain388_3(r, p) }
func chain388_3(r Router, p string) { deep_0(r, p) }

type m389 struct{}

func (m389) Register(r Router) { reg389(r, "/m389") }
func reg389(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain389_0(r, p)
}
func chain389_0(r Router, p string) { r.Get(p+"/deep0", h); chain389_1(r, p) }
func chain389_1(r Router, p string) { r.Get(p+"/deep1", h); chain389_2(r, p) }
func chain389_2(r Router, p string) { r.Get(p+"/deep2", h); chain389_3(r, p) }
func chain389_3(r Router, p string) { deep_0(r, p) }

type m390 struct{}

func (m390) Register(r Router) { reg390(r, "/m390") }
func reg390(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain390_0(r, p)
}
func chain390_0(r Router, p string) { r.Get(p+"/deep0", h); chain390_1(r, p) }
func chain390_1(r Router, p string) { r.Get(p+"/deep1", h); chain390_2(r, p) }
func chain390_2(r Router, p string) { r.Get(p+"/deep2", h); chain390_3(r, p) }
func chain390_3(r Router, p string) { deep_0(r, p) }

type m391 struct{}

func (m391) Register(r Router) { reg391(r, "/m391") }
func reg391(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain391_0(r, p)
}
func chain391_0(r Router, p string) { r.Get(p+"/deep0", h); chain391_1(r, p) }
func chain391_1(r Router, p string) { r.Get(p+"/deep1", h); chain391_2(r, p) }
func chain391_2(r Router, p string) { r.Get(p+"/deep2", h); chain391_3(r, p) }
func chain391_3(r Router, p string) { deep_0(r, p) }

type m392 struct{}

func (m392) Register(r Router) { reg392(r, "/m392") }
func reg392(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain392_0(r, p)
}
func chain392_0(r Router, p string) { r.Get(p+"/deep0", h); chain392_1(r, p) }
func chain392_1(r Router, p string) { r.Get(p+"/deep1", h); chain392_2(r, p) }
func chain392_2(r Router, p string) { r.Get(p+"/deep2", h); chain392_3(r, p) }
func chain392_3(r Router, p string) { deep_0(r, p) }

type m393 struct{}

func (m393) Register(r Router) { reg393(r, "/m393") }
func reg393(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain393_0(r, p)
}
func chain393_0(r Router, p string) { r.Get(p+"/deep0", h); chain393_1(r, p) }
func chain393_1(r Router, p string) { r.Get(p+"/deep1", h); chain393_2(r, p) }
func chain393_2(r Router, p string) { r.Get(p+"/deep2", h); chain393_3(r, p) }
func chain393_3(r Router, p string) { deep_0(r, p) }

type m394 struct{}

func (m394) Register(r Router) { reg394(r, "/m394") }
func reg394(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain394_0(r, p)
}
func chain394_0(r Router, p string) { r.Get(p+"/deep0", h); chain394_1(r, p) }
func chain394_1(r Router, p string) { r.Get(p+"/deep1", h); chain394_2(r, p) }
func chain394_2(r Router, p string) { r.Get(p+"/deep2", h); chain394_3(r, p) }
func chain394_3(r Router, p string) { deep_0(r, p) }

type m395 struct{}

func (m395) Register(r Router) { reg395(r, "/m395") }
func reg395(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain395_0(r, p)
}
func chain395_0(r Router, p string) { r.Get(p+"/deep0", h); chain395_1(r, p) }
func chain395_1(r Router, p string) { r.Get(p+"/deep1", h); chain395_2(r, p) }
func chain395_2(r Router, p string) { r.Get(p+"/deep2", h); chain395_3(r, p) }
func chain395_3(r Router, p string) { deep_0(r, p) }

type m396 struct{}

func (m396) Register(r Router) { reg396(r, "/m396") }
func reg396(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain396_0(r, p)
}
func chain396_0(r Router, p string) { r.Get(p+"/deep0", h); chain396_1(r, p) }
func chain396_1(r Router, p string) { r.Get(p+"/deep1", h); chain396_2(r, p) }
func chain396_2(r Router, p string) { r.Get(p+"/deep2", h); chain396_3(r, p) }
func chain396_3(r Router, p string) { deep_0(r, p) }

type m397 struct{}

func (m397) Register(r Router) { reg397(r, "/m397") }
func reg397(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain397_0(r, p)
}
func chain397_0(r Router, p string) { r.Get(p+"/deep0", h); chain397_1(r, p) }
func chain397_1(r Router, p string) { r.Get(p+"/deep1", h); chain397_2(r, p) }
func chain397_2(r Router, p string) { r.Get(p+"/deep2", h); chain397_3(r, p) }
func chain397_3(r Router, p string) { deep_0(r, p) }

type m398 struct{}

func (m398) Register(r Router) { reg398(r, "/m398") }
func reg398(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain398_0(r, p)
}
func chain398_0(r Router, p string) { r.Get(p+"/deep0", h); chain398_1(r, p) }
func chain398_1(r Router, p string) { r.Get(p+"/deep1", h); chain398_2(r, p) }
func chain398_2(r Router, p string) { r.Get(p+"/deep2", h); chain398_3(r, p) }
func chain398_3(r Router, p string) { deep_0(r, p) }

type m399 struct{}

func (m399) Register(r Router) { reg399(r, "/m399") }
func reg399(r Router, p string) {
	r.Get(p+"/items", h)
	r.Delete(p+"/items/{id}", h)
	r.Post(p+"/items", h)
	chain399_0(r, p)
}
func chain399_0(r Router, p string) { r.Get(p+"/deep0", h); chain399_1(r, p) }
func chain399_1(r Router, p string) { r.Get(p+"/deep1", h); chain399_2(r, p) }
func chain399_2(r Router, p string) { r.Get(p+"/deep2", h); chain399_3(r, p) }
func chain399_3(r Router, p string) { deep_0(r, p) }
func deep_0(r Router, p string)     { r.Get(p+"/d0", h); deep_1(r, p) }
func deep_1(r Router, p string)     { r.Get(p+"/d1", h); deep_2(r, p) }
func deep_2(r Router, p string)     { r.Get(p+"/d2", h); deep_3(r, p) }
func deep_3(r Router, p string)     { r.Get(p+"/d3", h); deep_4(r, p) }
func deep_4(r Router, p string)     { r.Get(p+"/d4", h); deep_5(r, p) }
func deep_5(r Router, p string)     { r.Get(p+"/d5", h); deep_6(r, p) }
func deep_6(r Router, p string)     { r.Get(p+"/d6", h); deep_7(r, p) }
func deep_7(r Router, p string)     { r.Get(p+"/d7", h); deep_8(r, p) }
func deep_8(r Router, p string)     { r.Get(p+"/d8", h); deep_9(r, p) }
func deep_9(r Router, p string)     { r.Get(p+"/d9", h); deep_10(r, p) }
func deep_10(r Router, p string)    { r.Get(p+"/d10", h); deep_11(r, p) }
func deep_11(r Router, p string)    { r.Get(p+"/d11", h); deep_12(r, p) }
func deep_12(r Router, p string)    { r.Get(p+"/d12", h); deep_13(r, p) }
func deep_13(r Router, p string)    { r.Get(p+"/d13", h); deep_14(r, p) }
func deep_14(r Router, p string)    { r.Get(p+"/d14", h); deep_15(r, p) }
func deep_15(r Router, p string)    { r.Get(p+"/d15", h); deep_16(r, p) }
func deep_16(r Router, p string)    { r.Get(p+"/d16", h); deep_17(r, p) }
func deep_17(r Router, p string)    { r.Get(p+"/d17", h); deep_18(r, p) }
func deep_18(r Router, p string)    { r.Get(p+"/d18", h); deep_19(r, p) }
func deep_19(r Router, p string)    { r.Get(p+"/d19", h); deep_20(r, p) }
func deep_20(r Router, p string)    { r.Get(p+"/d20", h); deep_21(r, p) }
func deep_21(r Router, p string)    { r.Get(p+"/d21", h); deep_22(r, p) }
func deep_22(r Router, p string)    { r.Get(p+"/d22", h); deep_23(r, p) }
func deep_23(r Router, p string)    { r.Get(p+"/d23", h); deep_24(r, p) }
func deep_24(r Router, p string)    { r.Get(p+"/d24", h); deep_25(r, p) }
func deep_25(r Router, p string)    { r.Get(p+"/d25", h); deep_26(r, p) }
func deep_26(r Router, p string)    { r.Get(p+"/d26", h); deep_27(r, p) }
func deep_27(r Router, p string)    { r.Get(p+"/d27", h); deep_28(r, p) }
func deep_28(r Router, p string)    { r.Get(p+"/d28", h); deep_29(r, p) }
func deep_29(r Router, p string)    { r.Get(p+"/d29", h); deep_30(r, p) }
func deep_30(r Router, p string)    { r.Get(p+"/d30", h); deep_31(r, p) }
func deep_31(r Router, p string)    { r.Get(p+"/d31", h); deep_32(r, p) }
func deep_32(r Router, p string)    { r.Get(p+"/d32", h); deep_33(r, p) }
func deep_33(r Router, p string)    { r.Get(p+"/d33", h); deep_34(r, p) }
func deep_34(r Router, p string)    { r.Get(p+"/d34", h); deep_35(r, p) }
func deep_35(r Router, p string)    { r.Get(p+"/d35", h); deep_36(r, p) }
func deep_36(r Router, p string)    { r.Get(p+"/d36", h); deep_37(r, p) }
func deep_37(r Router, p string)    { r.Get(p+"/d37", h); deep_38(r, p) }
func deep_38(r Router, p string)    { r.Get(p+"/d38", h); deep_39(r, p) }
func deep_39(r Router, p string)    { r.Get(p+"/d39", h); deep_40(r, p) }
func deep_40(r Router, p string)    {}
func all() []Module {
	return []Module{m0{}, m1{}, m2{}, m3{}, m4{}, m5{}, m6{}, m7{}, m8{}, m9{}, m10{}, m11{}, m12{}, m13{}, m14{}, m15{}, m16{}, m17{}, m18{}, m19{}, m20{}, m21{}, m22{}, m23{}, m24{}, m25{}, m26{}, m27{}, m28{}, m29{}, m30{}, m31{}, m32{}, m33{}, m34{}, m35{}, m36{}, m37{}, m38{}, m39{}, m40{}, m41{}, m42{}, m43{}, m44{}, m45{}, m46{}, m47{}, m48{}, m49{}, m50{}, m51{}, m52{}, m53{}, m54{}, m55{}, m56{}, m57{}, m58{}, m59{}, m60{}, m61{}, m62{}, m63{}, m64{}, m65{}, m66{}, m67{}, m68{}, m69{}, m70{}, m71{}, m72{}, m73{}, m74{}, m75{}, m76{}, m77{}, m78{}, m79{}, m80{}, m81{}, m82{}, m83{}, m84{}, m85{}, m86{}, m87{}, m88{}, m89{}, m90{}, m91{}, m92{}, m93{}, m94{}, m95{}, m96{}, m97{}, m98{}, m99{}, m100{}, m101{}, m102{}, m103{}, m104{}, m105{}, m106{}, m107{}, m108{}, m109{}, m110{}, m111{}, m112{}, m113{}, m114{}, m115{}, m116{}, m117{}, m118{}, m119{}, m120{}, m121{}, m122{}, m123{}, m124{}, m125{}, m126{}, m127{}, m128{}, m129{}, m130{}, m131{}, m132{}, m133{}, m134{}, m135{}, m136{}, m137{}, m138{}, m139{}, m140{}, m141{}, m142{}, m143{}, m144{}, m145{}, m146{}, m147{}, m148{}, m149{}, m150{}, m151{}, m152{}, m153{}, m154{}, m155{}, m156{}, m157{}, m158{}, m159{}, m160{}, m161{}, m162{}, m163{}, m164{}, m165{}, m166{}, m167{}, m168{}, m169{}, m170{}, m171{}, m172{}, m173{}, m174{}, m175{}, m176{}, m177{}, m178{}, m179{}, m180{}, m181{}, m182{}, m183{}, m184{}, m185{}, m186{}, m187{}, m188{}, m189{}, m190{}, m191{}, m192{}, m193{}, m194{}, m195{}, m196{}, m197{}, m198{}, m199{}, m200{}, m201{}, m202{}, m203{}, m204{}, m205{}, m206{}, m207{}, m208{}, m209{}, m210{}, m211{}, m212{}, m213{}, m214{}, m215{}, m216{}, m217{}, m218{}, m219{}, m220{}, m221{}, m222{}, m223{}, m224{}, m225{}, m226{}, m227{}, m228{}, m229{}, m230{}, m231{}, m232{}, m233{}, m234{}, m235{}, m236{}, m237{}, m238{}, m239{}, m240{}, m241{}, m242{}, m243{}, m244{}, m245{}, m246{}, m247{}, m248{}, m249{}, m250{}, m251{}, m252{}, m253{}, m254{}, m255{}, m256{}, m257{}, m258{}, m259{}, m260{}, m261{}, m262{}, m263{}, m264{}, m265{}, m266{}, m267{}, m268{}, m269{}, m270{}, m271{}, m272{}, m273{}, m274{}, m275{}, m276{}, m277{}, m278{}, m279{}, m280{}, m281{}, m282{}, m283{}, m284{}, m285{}, m286{}, m287{}, m288{}, m289{}, m290{}, m291{}, m292{}, m293{}, m294{}, m295{}, m296{}, m297{}, m298{}, m299{}, m300{}, m301{}, m302{}, m303{}, m304{}, m305{}, m306{}, m307{}, m308{}, m309{}, m310{}, m311{}, m312{}, m313{}, m314{}, m315{}, m316{}, m317{}, m318{}, m319{}, m320{}, m321{}, m322{}, m323{}, m324{}, m325{}, m326{}, m327{}, m328{}, m329{}, m330{}, m331{}, m332{}, m333{}, m334{}, m335{}, m336{}, m337{}, m338{}, m339{}, m340{}, m341{}, m342{}, m343{}, m344{}, m345{}, m346{}, m347{}, m348{}, m349{}, m350{}, m351{}, m352{}, m353{}, m354{}, m355{}, m356{}, m357{}, m358{}, m359{}, m360{}, m361{}, m362{}, m363{}, m364{}, m365{}, m366{}, m367{}, m368{}, m369{}, m370{}, m371{}, m372{}, m373{}, m374{}, m375{}, m376{}, m377{}, m378{}, m379{}, m380{}, m381{}, m382{}, m383{}, m384{}, m385{}, m386{}, m387{}, m388{}, m389{}, m390{}, m391{}, m392{}, m393{}, m394{}, m395{}, m396{}, m397{}, m398{}}
}
func main() {
	r := chi.NewRouter()
	r.Group(func(g chi.Router) {
		g.Use(RequireAuth)
		for _, m := range all() {
			m.Register(g)
		}
	})
	m399{}.Register(r)
	http.ListenAndServe(":8080", r)
}
