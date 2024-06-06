package ios

import (
	"fmt"
	"github.com/mobilemindtec/go-io/either"
	"github.com/mobilemindtec/go-io/option"
	"github.com/mobilemindtec/go-io/result"
	"github.com/mobilemindtec/go-io/runtime"
	"github.com/mobilemindtec/go-io/state"
	"github.com/mobilemindtec/go-io/types"
	"github.com/mobilemindtec/go-io/util"
	"log"
	"reflect"
	"runtime/debug"
)

type IOAttempt[A any] struct {
	value          *result.Result[*option.Option[A]]
	prevEffect     types.IOEffect
	fnResultOption func() *result.Result[*option.Option[A]]
	fnResult       func() *result.Result[A]
	fnOption       func() *option.Option[A]
	fnError        func() (A, error)

	fnUint      func()
	fnStateUint func(*state.State)

	fnStateResultOption func(*state.State) *result.Result[*option.Option[A]]
	fnStateResult       func(*state.State) *result.Result[A]
	fnStateOption       func(*state.State) *option.Option[A]
	fnStateError        func(*state.State) (A, error)
	fnPureState         func(*state.State) A

	fnStateIOIfEmpty func(*state.State) types.IORunnable
	fnStateIO        func(*state.State) types.IORunnable

	fnEitherResult      func() *result.Result[A]
	fnEitherStateResult func(*state.State) *result.Result[A]

	fnEither      func() A
	fnEitherState func(*state.State) A

	fnAuto interface{}

	fnExec             func(A)
	fnExecState        func(A, *state.State)
	fnExecIfEmpty      func()
	fnExecIfEmptyState func(*state.State)

	fnFlow      func(A) *result.Result[A]
	fnFlowState func(A, *state.State) *result.Result[A]

	fnFlowOption      func(A) *result.Result[*option.Option[A]]
	fnFlowOptionState func(A, *state.State) *result.Result[*option.Option[A]]

	state     *state.State
	debug     bool
	debugInfo *types.IODebugInfo
}

func NewAttempt[A any](f func() *result.Result[A]) *IOAttempt[A] {
	return &IOAttempt[A]{fnResult: f}
}

func NewAttemptOfOption[A any](f func() *option.Option[A]) *IOAttempt[A] {
	return &IOAttempt[A]{fnOption: f}
}

func NewAttemptOfResultOption[A any](f func() *result.Result[*option.Option[A]]) *IOAttempt[A] {
	return &IOAttempt[A]{fnResultOption: f}
}

func NewAttemptState[A any](f func(*state.State) *result.Result[A]) *IOAttempt[A] {
	return &IOAttempt[A]{fnStateResult: f}
}

func NewAttemptRunIOIfEmpty[A any](f func(*state.State) types.IORunnable) *IOAttempt[A] {
	return &IOAttempt[A]{fnStateIOIfEmpty: f}
}

func NewAttemptRunIO[A any](f func(*state.State) types.IORunnable) *IOAttempt[A] {
	return &IOAttempt[A]{fnStateIO: f}
}

func NewAttemptStateOfOption[A any](f func(*state.State) *option.Option[A]) *IOAttempt[A] {
	return &IOAttempt[A]{fnStateOption: f}
}

func NewAttemptStateOfResultOption[A any](f func(*state.State) *result.Result[*option.Option[A]]) *IOAttempt[A] {
	return &IOAttempt[A]{fnStateResultOption: f}
}

func NewAttemptOfResultEither[A any](f func() *result.Result[A]) *IOAttempt[A] {
	return &IOAttempt[A]{fnEitherResult: f}
}

func NewAttemptStateOfResultEither[A any](f func(*state.State) *result.Result[A]) *IOAttempt[A] {
	return &IOAttempt[A]{fnEitherStateResult: f}
}

func NewAttemptOfEither[A any](f func() A) *IOAttempt[A] {
	return &IOAttempt[A]{fnEither: f}
}

func NewAttemptStateOfEither[A any](f func(*state.State) A) *IOAttempt[A] {
	return &IOAttempt[A]{fnEitherState: f}
}

func NewAttemptAuto[A any](f interface{}) *IOAttempt[A] {
	return &IOAttempt[A]{fnAuto: f}
}

func NewAttemptOfUnit[A any](f func()) *IOAttempt[A] {
	return &IOAttempt[A]{fnUint: f}
}

func NewAttemptStateOfUnit[A any](f func(*state.State)) *IOAttempt[A] {
	return &IOAttempt[A]{fnStateUint: f}
}

func NewAttemptOfError[A any](f func() (A, error)) *IOAttempt[A] {
	return &IOAttempt[A]{fnError: f}
}

func NewAttemptStateOfError[A any](f func(*state.State) (A, error)) *IOAttempt[A] {
	return &IOAttempt[A]{fnStateError: f}
}

func NewAttemptPureState[A any](f func(*state.State) A) *IOAttempt[A] {
	return &IOAttempt[A]{fnPureState: f}
}

func NewAttemptExec[A any](f func(A)) *IOAttempt[A] {
	return &IOAttempt[A]{fnExec: f}
}

func NewAttemptExecState[A any](f func(A, *state.State)) *IOAttempt[A] {
	return &IOAttempt[A]{fnExecState: f}
}

func NewAttemptExecIfEmpty[A any](f func()) *IOAttempt[A] {
	return &IOAttempt[A]{fnExecIfEmpty: f}
}

func NewAttemptExecIfEmptyState[A any](f func(*state.State)) *IOAttempt[A] {
	return &IOAttempt[A]{fnExecIfEmptyState: f}
}

func NewAttemptFlowOfResult[A any](f func(A) *result.Result[A]) *IOAttempt[A] {
	return &IOAttempt[A]{fnFlow: f}
}

func NewAttemptFlowStateOfResult[A any](f func(A, *state.State) *result.Result[A]) *IOAttempt[A] {
	return &IOAttempt[A]{fnFlowState: f}
}

func NewAttemptFlowOfResultOption[A any](f func(A) *result.Result[*option.Option[A]]) *IOAttempt[A] {
	return &IOAttempt[A]{fnFlowOption: f}
}

func NewAttemptFlowStateOfResultOption[A any](f func(A, *state.State) *result.Result[*option.Option[A]]) *IOAttempt[A] {
	return &IOAttempt[A]{fnFlowOptionState: f}
}

func (this *IOAttempt[T]) Lift() *types.IO[T] {
	return types.NewIO[T]().Effects(this)
}

func (this *IOAttempt[A]) SetDebug(b bool) {
	this.debug = b
}

func (this *IOAttempt[T]) SetDebugInfo(info *types.IODebugInfo) {
	this.debugInfo = info
}

func (this *IOAttempt[T]) GetDebugInfo() *types.IODebugInfo {
	return this.debugInfo
}

func (this *IOAttempt[T]) TypeIn() reflect.Type {
	return reflect.TypeFor[T]()
}
func (this *IOAttempt[T]) TypeOut() reflect.Type {
	return reflect.TypeFor[T]()
}

func (this *IOAttempt[A]) String() string {
	return fmt.Sprintf("Attempt(fn=%v, value=%v)", this.getFuncName(), this.value.String())
}

func (this *IOAttempt[A]) SetState(st *state.State) {
	this.state = st
}

func (this *IOAttempt[A]) SetPrevEffect(prev types.IOEffect) {
	this.prevEffect = prev
}

func (this *IOAttempt[A]) GetPrevEffect() *option.Option[types.IOEffect] {
	return option.Of(this.prevEffect)
}

func (this *IOAttempt[A]) GetResult() types.ResultOptionAny {
	return this.value.ToResultOfOption()
}

func (this *IOAttempt[A]) UnsafeRun() types.IOEffect {

	var currEff interface{} = this
	prevEff := this.GetPrevEffect()
	this.value = result.OfValue(option.None[A]())
	execute := true
	isEmpty := false
	hasPrev := prevEff.NonEmpty()

	if hasPrev {
		prev := prevEff.Get()
		if prev.GetResult().IsError() {
			this.value = result.OfError[*option.Option[A]](
				prevEff.Get().GetResult().Failure())
			execute = false
		} else {
			execute = prev.GetResult().Get().NonEmpty()
			isEmpty = prev.GetResult().Get().IsEmpty()
		}
	}

	defer func() {
		if r := recover(); r != nil {

			err := types.NewIOError(fmt.Sprintf("%v", r), debug.Stack())

			if this.debug {
				if this.debugInfo != nil {
					log.Printf("[DEBUG IOAttempt]=>> added in: %v:%v", this.debugInfo.Filename, this.debugInfo.Line)
				}
				log.Printf("[DEBUG IOAttempt]=>> Error: %v\n", err.Error())
				log.Printf("[DEBUG IOAttempt]=>> StackTrace: %v\n", err.StackTrace)
			}

			this.value = result.OfError[*option.Option[A]](err)
		}
	}()

	execOnEmpty := this.fnExecIfEmpty != nil ||
		this.fnExecIfEmptyState != nil ||
		this.fnStateIOIfEmpty != nil

	if isEmpty { // not error
		if execOnEmpty {
			if this.fnExecIfEmpty != nil {
				this.fnExecIfEmpty()
			} else if this.fnExecIfEmptyState != nil {
				this.fnExecIfEmptyState(this.state)
			} else if this.fnStateIOIfEmpty != nil {

				if this.debug {
					log.Printf("IOAttempt >> RunIO IfEmpty\n")
				}

				runnableIO := this.fnStateIOIfEmpty(this.state)

				if this.debug {
					runnableIO.SetDebug(this.debug)
				}

				this.value = runtime.New[A](runnableIO).UnsafeRun()

			} else {
				panic("func not found (empty)")
			}
		}
	} else {

		// if it has an "if empty" exec and last result is not empty,
		// so preserves the last result
		if execOnEmpty {
			r := prevEff.Get().GetResult()
			val := r.Get().GetValue()

			if effValue, ok := val.(A); ok {
				this.value = result.OfValue(option.Some(effValue))
			} else {
				util.PanicCastType("IOAttempt",
					reflect.TypeOf(val), reflect.TypeFor[A]())
			}
		}
	}

	if execute {

		if this.fnStateIO != nil {

			if this.debug {
				log.Printf("IOAttempt >> RunIO\n")
			}

			runnableIO := this.fnStateIO(this.state)

			if this.debug {
				runnableIO.SetDebug(this.debug)
			}

			this.value = runtime.New[A](runnableIO).UnsafeRun()

		} else if this.fnExec != nil || this.fnExecState != nil && hasPrev {

			r := prevEff.Get().GetResult()
			val := r.Get().GetValue()

			if effValue, ok := val.(A); ok {
				if this.fnExec != nil {
					this.fnExec(effValue)
				} else {
					this.fnExecState(effValue, this.state)
				}
			} else {
				util.PanicCastType("IOAttempt",
					reflect.TypeOf(val), reflect.TypeFor[A]())

			}

		} else if this.fnFlow != nil || this.fnFlowState != nil && hasPrev {

			r := prevEff.Get().GetResult()
			val := r.Get().GetValue()

			if effValue, ok := val.(A); ok {
				var res *result.Result[A]
				if this.fnFlow != nil {
					res = this.fnFlow(effValue)
				} else {
					res = this.fnFlowState(effValue, this.state)
				}
				if res.IsError() {
					this.value = result.OfError[*option.Option[A]](res.GetError())
				} else {
					this.value = result.OfValue(option.Of(res.Get()))
				}
			} else {
				util.PanicCastType("IOAttempt",
					reflect.TypeOf(val), reflect.TypeFor[A]())

			}

		} else if this.fnFlowOption != nil || this.fnFlowOptionState != nil && hasPrev {

			r := prevEff.Get().GetResult()
			val := r.Get().GetValue()

			if effValue, ok := val.(A); ok {
				if this.fnFlowOption != nil {
					this.value = this.fnFlowOption(effValue)
				} else {
					this.value = this.fnFlowOptionState(effValue, this.state)
				}
			} else {
				util.PanicCastType("IOAttempt",
					reflect.TypeOf(val), reflect.TypeFor[A]())

			}

		} else if this.fnResultOption != nil {
			this.value = this.fnResultOption()
		} else if this.fnOption != nil {
			this.value = result.OfValue(this.fnOption())
		} else if this.fnResult != nil {
			r := this.fnResult()
			if r.HasError() {
				this.value = result.OfError[*option.Option[A]](r.GetError())
			} else {
				this.value = result.OfValue(option.Some(r.Get()))
			}
		} else if this.fnError != nil {
			this.value = result.TryOption(this.fnError)
		} else if this.fnStateResultOption != nil {
			this.value = this.fnStateResultOption(this.state)
		} else if this.fnStateOption != nil {
			this.value = result.OfValue(this.fnStateOption(this.state))
		} else if this.fnStateResult != nil {
			r := this.fnStateResult(this.state)
			if r.HasError() {
				this.value = result.OfError[*option.Option[A]](r.GetError())
			} else {
				this.value = result.OfValue(option.Some(r.Get()))
			}
		} else if this.fnStateError != nil {
			this.value = result.TryOption(
				func() (A, error) {
					return this.fnStateError(this.state)
				})
		} else if this.fnPureState != nil {
			this.value = result.OfValue(option.Of(this.fnPureState(this.state)))
		} else if this.fnUint != nil {
			this.fnUint()
			var unit interface{} = types.OfUnit()
			this.value = result.OfValue(option.Some(unit.(A)))
		} else if this.fnStateUint != nil {
			this.fnStateUint(this.state)
			var unit interface{} = types.OfUnit()
			this.value = result.OfValue(option.Some(unit.(A)))
		} else if this.fnEitherResult != nil || this.fnEitherStateResult != nil { // either
			var ret *result.Result[A]
			if this.fnEitherResult != nil {
				ret = this.fnEitherResult()
			} else {
				ret = this.fnEitherStateResult(this.state)
			}
			if ret.HasError() {
				this.value = result.OfError[*option.Option[A]](ret.Failure())
			} else {
				val := ret.GetValue()
				if _, ok := val.(either.IEither); ok {
					if eia, ok := val.(A); ok {
						this.value = result.OfValue[*option.Option[A]](option.Some(eia))
					} else {
						util.PanicCastType("IOAttempt",
							reflect.TypeOf(val), reflect.TypeFor[A]())
					}
				} else {
					util.PanicCastType("IOAttempt",
						reflect.TypeOf(val), reflect.TypeFor[either.IEither]())
				}
			}
		} else if this.fnEither != nil || this.fnEitherState != nil { // either
			var ret interface{}
			if this.fnEither != nil {
				ret = this.fnEither()
			} else {
				ret = this.fnEitherState(this.state)
			}
			if _, ok := ret.(either.IEither); ok {
				if eia, ok := ret.(A); ok {
					this.value = result.OfValue[*option.Option[A]](option.Some(eia))
				} else {
					util.PanicCastType("IOAttempt",
						reflect.TypeOf(ret), reflect.TypeFor[A]())
				}
			} else {
				util.PanicCastType("IOAttempt",
					reflect.TypeOf(ret), reflect.TypeFor[either.IEither]())
			}
		} else if this.fnAuto != nil {

			info := util.NewFuncInfo(this.fnAuto)
			var fnParams []reflect.Value
			stateCopy := this.state.Copy()
			for i := 0; i < info.ArgsCount; i++ {
				_, val := state.LookupVar(stateCopy, info.ArgType(i), true)
				fnParams = append(fnParams, val)
			}

			fnResults := info.Call(fnParams)

			switch len(fnResults) {
			case 0:
				var unit interface{} = types.OfUnit()
				this.value = result.OfValue(option.Some(unit.(A)))
			case 1: // should be a Result[A]

				rVal := fnResults[0]

				if util.CanNil(rVal.Kind()) && rVal.IsNil() {
					panic("func result can't be null")
				}

				retValue := rVal.Interface()

				if r, ok := retValue.(result.IResult); ok {

					if r.HasError() {
						this.value = result.OfError[*option.Option[A]](r.GetError())
					} else {
						if o, ok := r.GetValue().(option.IOption); ok {
							if o.IsEmpty() {
								this.value = result.OfValue(option.None[A]())
							} else {

								if val, ok := o.GetValue().(A); ok {
									this.value = result.OfValue(option.Some(val))
								} else {
									util.PanicCastType("IOAttempt",
										reflect.TypeOf(o.GetValue()), reflect.TypeFor[A]())
								}
							}
						} else if val, ok := r.GetValue().(A); ok {
							this.value = result.OfValue(option.Some(val))
						} else {
							panic("wrong result type")
						}
					}

				} else if opt, ok := retValue.(option.IOption); ok {

					if opt.IsEmpty() {
						this.value = result.OfValue(option.None[A]())
					} else {
						if val, ok := opt.GetValue().(A); ok {
							this.value = result.OfValue(option.Some(val))
						} else {
							util.PanicCastType("IOAttempt",
								reflect.TypeOf(opt.GetValue()), reflect.TypeFor[A]())
						}
					}

				} else if _, ok := retValue.(either.IEither); ok {
					if eia, ok := retValue.(A); ok {
						this.value = result.OfValue[*option.Option[A]](option.Some(eia))
					} else {
						util.PanicCastType("IOAttempt",
							reflect.TypeOf(retValue), reflect.TypeFor[A]())
					}
				} else {
					panic("wrong result type")
				}

			case 2: // should be a (A, error)

				rA := fnResults[0]
				rE := fnResults[1]
				canNil := util.CanNil(rA.Kind())

				if (canNil && rA.IsNil()) && rE.IsNil() {
					panic("func result can't be null")
				}

				retTypeA := rA.Interface()
				retTypeE := rE.Interface()

				var valA A
				var valE error
				var ok bool

				if !canNil || !rA.IsNil() {
					valA, ok = rA.Interface().(A)
					if !ok {
						panic(fmt.Sprintf("cat't cast %v to %b", retTypeA, reflect.TypeOf(valA)))
					}
				}

				if !rE.IsNil() {
					valE, ok = rE.Interface().(error)
					if !ok {
						panic(fmt.Sprintf("cat't cast %v to %b", retTypeE, reflect.TypeOf(valE)))
					}
				}

				this.value = result.TryOption(func() (A, error) {
					return valA, valE
				})

			default:
				panic("func should be result.Result[*option.Option[A]] or (A, error)")
			}
		}
	}

	if this.debug {
		log.Printf("%v\n", this.String())
	}

	return currEff.(types.IOEffect)
}

func (this *IOAttempt[A]) getFuncName() string {
	if this.fnResultOption != nil {
		return "fnResultOption"
	}
	if this.fnResult != nil {
		return "fnResult"
	}
	if this.fnOption != nil {
		return "fnOption"
	}
	if this.fnError != nil {
		return "fnError"
	}
	if this.fnUint != nil {
		return "fnUint"
	}
	if this.fnStateUint != nil {
		return "fnStateUint"
	}
	if this.fnStateResultOption != nil {
		return "fnStateResultOption"
	}
	if this.fnStateResult != nil {
		return "fnStateResult"
	}
	if this.fnStateOption != nil {
		return "fnStateOption"
	}
	if this.fnStateError != nil {
		return "fnStateError"
	}
	if this.fnPureState != nil {
		return "fnPureState"
	}
	if this.fnStateIOIfEmpty != nil {
		return "fnStateIOIfEmpty"
	}
	if this.fnStateIO != nil {
		return "fnStateIO"
	}
	if this.fnEitherResult != nil {
		return "fnEitherResult"
	}
	if this.fnEitherStateResult != nil {
		return "fnEitherStateResult"
	}
	if this.fnEither != nil {
		return "fnEither"
	}
	if this.fnEitherState != nil {
		return "fnEitherState"
	}
	if this.fnAuto != nil {
		return "fnAuto"
	}
	if this.fnExec != nil {
		return "fnExec"
	}
	if this.fnExecState != nil {
		return "fnExecState"
	}
	if this.fnExecIfEmpty != nil {
		return "fnExecIfEmpty"
	}
	if this.fnExecIfEmptyState != nil {
		return "fnExecIfEmptyState"
	}
	if this.fnFlow != nil {
		return "fnFlow"
	}
	if this.fnFlowState != nil {
		return "fnFlowState"
	}
	if this.fnFlowOption != nil {
		return "fnFlowOption"
	}
	if this.fnFlowOptionState != nil {
		return "fnFlowOptionState"
	}
	return "-"
}