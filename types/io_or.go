package types

import (
	"fmt"
	"github.com/mobilemindtec/go-io/option"
	"github.com/mobilemindtec/go-io/result"
	"github.com/mobilemindtec/go-io/state"
	"github.com/mobilemindtec/go-io/util"
	"log"
	"reflect"
)

type IOOr[A any] struct {
	value      *result.Result[*option.Option[A]]
	prevEffect IOEffect
	f          func() A
	debug      bool
	state      *state.State
}

func NewOr[A any](f func() A) *IOOr[A] {
	return &IOOr[A]{f: f}
}

func (this *IOOr[A]) SetState(st *state.State) {
	this.state = st
}

func (this *IOOr[A]) SetDebug(b bool) {
	this.debug = b
}

func (this *IOOr[A]) String() string {
	return fmt.Sprintf("Or(%v)", this.value.String())
}

func (this *IOOr[A]) SetPrevEffect(prev IOEffect) {
	this.prevEffect = prev
}

func (this *IOOr[A]) GetPrevEffect() *option.Option[IOEffect] {
	return option.Of(this.prevEffect)
}

func (this *IOOr[A]) GetResult() ResultOptionAny {
	return this.value.ToResultOfOption()
}

func (this *IOOr[A]) UnsafeRun() IOEffect {
	var currEff interface{} = this
	prevEff := this.GetPrevEffect()
	this.value = result.OfValue(option.None[A]())

	if prevEff.NonEmpty() {
		r := prevEff.Get().GetResult()
		if r.IsError() {
			this.value = result.OfError[*option.Option[A]](r.Failure())
		} else if r.Get().Empty() {
			this.value = result.OfValue(option.Some(this.f()))
		} else {
			val := r.Get().GetValue()
			if effValue, ok := val.(A); ok {
				this.value = result.OfValue(option.Some(effValue))
			} else {
				util.PanicCastType("IOOr",
					reflect.TypeOf(val), reflect.TypeFor[A]())

			}
		}
	}

	if this.debug {
		log.Printf("%v\n", this.String())
	}
	return currEff.(IOEffect)
}
