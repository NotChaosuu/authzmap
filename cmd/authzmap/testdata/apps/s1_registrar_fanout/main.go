package main

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

func h(w http.ResponseWriter, r *http.Request) {}
func reg0(r chi.Router, d int) {
	r.Get("/p0", h)
	if d > 0 {
		reg1(r.With(mw0), d-1)
		reg2(r, d-1)
		r.Route("/n0", func(r chi.Router) { reg3(r, d-1) })
	}
}
func reg1(r chi.Router, d int) {
	r.Get("/p1", h)
	if d > 0 {
		reg2(r.With(mw1), d-1)
		reg3(r, d-1)
		r.Route("/n1", func(r chi.Router) { reg10(r, d-1) })
	}
}
func reg2(r chi.Router, d int) {
	r.Get("/p2", h)
	if d > 0 {
		reg3(r.With(mw2), d-1)
		reg4(r, d-1)
		r.Route("/n2", func(r chi.Router) { reg17(r, d-1) })
	}
}
func reg3(r chi.Router, d int) {
	r.Get("/p3", h)
	if d > 0 {
		reg4(r.With(mw3), d-1)
		reg5(r, d-1)
		r.Route("/n3", func(r chi.Router) { reg24(r, d-1) })
	}
}
func reg4(r chi.Router, d int) {
	r.Get("/p4", h)
	if d > 0 {
		reg5(r.With(mw4), d-1)
		reg6(r, d-1)
		r.Route("/n4", func(r chi.Router) { reg31(r, d-1) })
	}
}
func reg5(r chi.Router, d int) {
	r.Get("/p5", h)
	if d > 0 {
		reg6(r.With(mw0), d-1)
		reg7(r, d-1)
		r.Route("/n5", func(r chi.Router) { reg38(r, d-1) })
	}
}
func reg6(r chi.Router, d int) {
	r.Get("/p6", h)
	if d > 0 {
		reg7(r.With(mw1), d-1)
		reg8(r, d-1)
		r.Route("/n6", func(r chi.Router) { reg45(r, d-1) })
	}
}
func reg7(r chi.Router, d int) {
	r.Get("/p7", h)
	if d > 0 {
		reg8(r.With(mw2), d-1)
		reg9(r, d-1)
		r.Route("/n7", func(r chi.Router) { reg52(r, d-1) })
	}
}
func reg8(r chi.Router, d int) {
	r.Get("/p8", h)
	if d > 0 {
		reg9(r.With(mw3), d-1)
		reg10(r, d-1)
		r.Route("/n8", func(r chi.Router) { reg59(r, d-1) })
	}
}
func reg9(r chi.Router, d int) {
	r.Get("/p9", h)
	if d > 0 {
		reg10(r.With(mw4), d-1)
		reg11(r, d-1)
		r.Route("/n9", func(r chi.Router) { reg6(r, d-1) })
	}
}
func reg10(r chi.Router, d int) {
	r.Get("/p10", h)
	if d > 0 {
		reg11(r.With(mw0), d-1)
		reg12(r, d-1)
		r.Route("/n10", func(r chi.Router) { reg13(r, d-1) })
	}
}
func reg11(r chi.Router, d int) {
	r.Get("/p11", h)
	if d > 0 {
		reg12(r.With(mw1), d-1)
		reg13(r, d-1)
		r.Route("/n11", func(r chi.Router) { reg20(r, d-1) })
	}
}
func reg12(r chi.Router, d int) {
	r.Get("/p12", h)
	if d > 0 {
		reg13(r.With(mw2), d-1)
		reg14(r, d-1)
		r.Route("/n12", func(r chi.Router) { reg27(r, d-1) })
	}
}
func reg13(r chi.Router, d int) {
	r.Get("/p13", h)
	if d > 0 {
		reg14(r.With(mw3), d-1)
		reg15(r, d-1)
		r.Route("/n13", func(r chi.Router) { reg34(r, d-1) })
	}
}
func reg14(r chi.Router, d int) {
	r.Get("/p14", h)
	if d > 0 {
		reg15(r.With(mw4), d-1)
		reg16(r, d-1)
		r.Route("/n14", func(r chi.Router) { reg41(r, d-1) })
	}
}
func reg15(r chi.Router, d int) {
	r.Get("/p15", h)
	if d > 0 {
		reg16(r.With(mw0), d-1)
		reg17(r, d-1)
		r.Route("/n15", func(r chi.Router) { reg48(r, d-1) })
	}
}
func reg16(r chi.Router, d int) {
	r.Get("/p16", h)
	if d > 0 {
		reg17(r.With(mw1), d-1)
		reg18(r, d-1)
		r.Route("/n16", func(r chi.Router) { reg55(r, d-1) })
	}
}
func reg17(r chi.Router, d int) {
	r.Get("/p17", h)
	if d > 0 {
		reg18(r.With(mw2), d-1)
		reg19(r, d-1)
		r.Route("/n17", func(r chi.Router) { reg2(r, d-1) })
	}
}
func reg18(r chi.Router, d int) {
	r.Get("/p18", h)
	if d > 0 {
		reg19(r.With(mw3), d-1)
		reg20(r, d-1)
		r.Route("/n18", func(r chi.Router) { reg9(r, d-1) })
	}
}
func reg19(r chi.Router, d int) {
	r.Get("/p19", h)
	if d > 0 {
		reg20(r.With(mw4), d-1)
		reg21(r, d-1)
		r.Route("/n19", func(r chi.Router) { reg16(r, d-1) })
	}
}
func reg20(r chi.Router, d int) {
	r.Get("/p20", h)
	if d > 0 {
		reg21(r.With(mw0), d-1)
		reg22(r, d-1)
		r.Route("/n20", func(r chi.Router) { reg23(r, d-1) })
	}
}
func reg21(r chi.Router, d int) {
	r.Get("/p21", h)
	if d > 0 {
		reg22(r.With(mw1), d-1)
		reg23(r, d-1)
		r.Route("/n21", func(r chi.Router) { reg30(r, d-1) })
	}
}
func reg22(r chi.Router, d int) {
	r.Get("/p22", h)
	if d > 0 {
		reg23(r.With(mw2), d-1)
		reg24(r, d-1)
		r.Route("/n22", func(r chi.Router) { reg37(r, d-1) })
	}
}
func reg23(r chi.Router, d int) {
	r.Get("/p23", h)
	if d > 0 {
		reg24(r.With(mw3), d-1)
		reg25(r, d-1)
		r.Route("/n23", func(r chi.Router) { reg44(r, d-1) })
	}
}
func reg24(r chi.Router, d int) {
	r.Get("/p24", h)
	if d > 0 {
		reg25(r.With(mw4), d-1)
		reg26(r, d-1)
		r.Route("/n24", func(r chi.Router) { reg51(r, d-1) })
	}
}
func reg25(r chi.Router, d int) {
	r.Get("/p25", h)
	if d > 0 {
		reg26(r.With(mw0), d-1)
		reg27(r, d-1)
		r.Route("/n25", func(r chi.Router) { reg58(r, d-1) })
	}
}
func reg26(r chi.Router, d int) {
	r.Get("/p26", h)
	if d > 0 {
		reg27(r.With(mw1), d-1)
		reg28(r, d-1)
		r.Route("/n26", func(r chi.Router) { reg5(r, d-1) })
	}
}
func reg27(r chi.Router, d int) {
	r.Get("/p27", h)
	if d > 0 {
		reg28(r.With(mw2), d-1)
		reg29(r, d-1)
		r.Route("/n27", func(r chi.Router) { reg12(r, d-1) })
	}
}
func reg28(r chi.Router, d int) {
	r.Get("/p28", h)
	if d > 0 {
		reg29(r.With(mw3), d-1)
		reg30(r, d-1)
		r.Route("/n28", func(r chi.Router) { reg19(r, d-1) })
	}
}
func reg29(r chi.Router, d int) {
	r.Get("/p29", h)
	if d > 0 {
		reg30(r.With(mw4), d-1)
		reg31(r, d-1)
		r.Route("/n29", func(r chi.Router) { reg26(r, d-1) })
	}
}
func reg30(r chi.Router, d int) {
	r.Get("/p30", h)
	if d > 0 {
		reg31(r.With(mw0), d-1)
		reg32(r, d-1)
		r.Route("/n30", func(r chi.Router) { reg33(r, d-1) })
	}
}
func reg31(r chi.Router, d int) {
	r.Get("/p31", h)
	if d > 0 {
		reg32(r.With(mw1), d-1)
		reg33(r, d-1)
		r.Route("/n31", func(r chi.Router) { reg40(r, d-1) })
	}
}
func reg32(r chi.Router, d int) {
	r.Get("/p32", h)
	if d > 0 {
		reg33(r.With(mw2), d-1)
		reg34(r, d-1)
		r.Route("/n32", func(r chi.Router) { reg47(r, d-1) })
	}
}
func reg33(r chi.Router, d int) {
	r.Get("/p33", h)
	if d > 0 {
		reg34(r.With(mw3), d-1)
		reg35(r, d-1)
		r.Route("/n33", func(r chi.Router) { reg54(r, d-1) })
	}
}
func reg34(r chi.Router, d int) {
	r.Get("/p34", h)
	if d > 0 {
		reg35(r.With(mw4), d-1)
		reg36(r, d-1)
		r.Route("/n34", func(r chi.Router) { reg1(r, d-1) })
	}
}
func reg35(r chi.Router, d int) {
	r.Get("/p35", h)
	if d > 0 {
		reg36(r.With(mw0), d-1)
		reg37(r, d-1)
		r.Route("/n35", func(r chi.Router) { reg8(r, d-1) })
	}
}
func reg36(r chi.Router, d int) {
	r.Get("/p36", h)
	if d > 0 {
		reg37(r.With(mw1), d-1)
		reg38(r, d-1)
		r.Route("/n36", func(r chi.Router) { reg15(r, d-1) })
	}
}
func reg37(r chi.Router, d int) {
	r.Get("/p37", h)
	if d > 0 {
		reg38(r.With(mw2), d-1)
		reg39(r, d-1)
		r.Route("/n37", func(r chi.Router) { reg22(r, d-1) })
	}
}
func reg38(r chi.Router, d int) {
	r.Get("/p38", h)
	if d > 0 {
		reg39(r.With(mw3), d-1)
		reg40(r, d-1)
		r.Route("/n38", func(r chi.Router) { reg29(r, d-1) })
	}
}
func reg39(r chi.Router, d int) {
	r.Get("/p39", h)
	if d > 0 {
		reg40(r.With(mw4), d-1)
		reg41(r, d-1)
		r.Route("/n39", func(r chi.Router) { reg36(r, d-1) })
	}
}
func reg40(r chi.Router, d int) {
	r.Get("/p40", h)
	if d > 0 {
		reg41(r.With(mw0), d-1)
		reg42(r, d-1)
		r.Route("/n40", func(r chi.Router) { reg43(r, d-1) })
	}
}
func reg41(r chi.Router, d int) {
	r.Get("/p41", h)
	if d > 0 {
		reg42(r.With(mw1), d-1)
		reg43(r, d-1)
		r.Route("/n41", func(r chi.Router) { reg50(r, d-1) })
	}
}
func reg42(r chi.Router, d int) {
	r.Get("/p42", h)
	if d > 0 {
		reg43(r.With(mw2), d-1)
		reg44(r, d-1)
		r.Route("/n42", func(r chi.Router) { reg57(r, d-1) })
	}
}
func reg43(r chi.Router, d int) {
	r.Get("/p43", h)
	if d > 0 {
		reg44(r.With(mw3), d-1)
		reg45(r, d-1)
		r.Route("/n43", func(r chi.Router) { reg4(r, d-1) })
	}
}
func reg44(r chi.Router, d int) {
	r.Get("/p44", h)
	if d > 0 {
		reg45(r.With(mw4), d-1)
		reg46(r, d-1)
		r.Route("/n44", func(r chi.Router) { reg11(r, d-1) })
	}
}
func reg45(r chi.Router, d int) {
	r.Get("/p45", h)
	if d > 0 {
		reg46(r.With(mw0), d-1)
		reg47(r, d-1)
		r.Route("/n45", func(r chi.Router) { reg18(r, d-1) })
	}
}
func reg46(r chi.Router, d int) {
	r.Get("/p46", h)
	if d > 0 {
		reg47(r.With(mw1), d-1)
		reg48(r, d-1)
		r.Route("/n46", func(r chi.Router) { reg25(r, d-1) })
	}
}
func reg47(r chi.Router, d int) {
	r.Get("/p47", h)
	if d > 0 {
		reg48(r.With(mw2), d-1)
		reg49(r, d-1)
		r.Route("/n47", func(r chi.Router) { reg32(r, d-1) })
	}
}
func reg48(r chi.Router, d int) {
	r.Get("/p48", h)
	if d > 0 {
		reg49(r.With(mw3), d-1)
		reg50(r, d-1)
		r.Route("/n48", func(r chi.Router) { reg39(r, d-1) })
	}
}
func reg49(r chi.Router, d int) {
	r.Get("/p49", h)
	if d > 0 {
		reg50(r.With(mw4), d-1)
		reg51(r, d-1)
		r.Route("/n49", func(r chi.Router) { reg46(r, d-1) })
	}
}
func reg50(r chi.Router, d int) {
	r.Get("/p50", h)
	if d > 0 {
		reg51(r.With(mw0), d-1)
		reg52(r, d-1)
		r.Route("/n50", func(r chi.Router) { reg53(r, d-1) })
	}
}
func reg51(r chi.Router, d int) {
	r.Get("/p51", h)
	if d > 0 {
		reg52(r.With(mw1), d-1)
		reg53(r, d-1)
		r.Route("/n51", func(r chi.Router) { reg0(r, d-1) })
	}
}
func reg52(r chi.Router, d int) {
	r.Get("/p52", h)
	if d > 0 {
		reg53(r.With(mw2), d-1)
		reg54(r, d-1)
		r.Route("/n52", func(r chi.Router) { reg7(r, d-1) })
	}
}
func reg53(r chi.Router, d int) {
	r.Get("/p53", h)
	if d > 0 {
		reg54(r.With(mw3), d-1)
		reg55(r, d-1)
		r.Route("/n53", func(r chi.Router) { reg14(r, d-1) })
	}
}
func reg54(r chi.Router, d int) {
	r.Get("/p54", h)
	if d > 0 {
		reg55(r.With(mw4), d-1)
		reg56(r, d-1)
		r.Route("/n54", func(r chi.Router) { reg21(r, d-1) })
	}
}
func reg55(r chi.Router, d int) {
	r.Get("/p55", h)
	if d > 0 {
		reg56(r.With(mw0), d-1)
		reg57(r, d-1)
		r.Route("/n55", func(r chi.Router) { reg28(r, d-1) })
	}
}
func reg56(r chi.Router, d int) {
	r.Get("/p56", h)
	if d > 0 {
		reg57(r.With(mw1), d-1)
		reg58(r, d-1)
		r.Route("/n56", func(r chi.Router) { reg35(r, d-1) })
	}
}
func reg57(r chi.Router, d int) {
	r.Get("/p57", h)
	if d > 0 {
		reg58(r.With(mw2), d-1)
		reg59(r, d-1)
		r.Route("/n57", func(r chi.Router) { reg42(r, d-1) })
	}
}
func reg58(r chi.Router, d int) {
	r.Get("/p58", h)
	if d > 0 {
		reg59(r.With(mw3), d-1)
		reg0(r, d-1)
		r.Route("/n58", func(r chi.Router) { reg49(r, d-1) })
	}
}
func reg59(r chi.Router, d int) {
	r.Get("/p59", h)
	if d > 0 {
		reg0(r.With(mw4), d-1)
		reg1(r, d-1)
		r.Route("/n59", func(r chi.Router) { reg56(r, d-1) })
	}
}
func mw0(n http.Handler) http.Handler { return n }
func mw1(n http.Handler) http.Handler { return n }
func mw2(n http.Handler) http.Handler { return n }
func mw3(n http.Handler) http.Handler { return n }
func mw4(n http.Handler) http.Handler { return n }
func factory(r chi.Router) chi.Router {
	s := chi.NewRouter()
	r.Mount("/f", s)
	s.Mount("/g", factory(s))
	return s
}
func main() { r := chi.NewRouter(); reg0(r, 50); factory(r); http.ListenAndServe(":1", r) }
