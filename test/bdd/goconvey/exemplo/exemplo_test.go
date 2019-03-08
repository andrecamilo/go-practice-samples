package exemplo

import (
	"strconv"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSingleScope(t *testing.T) {
	output := prepare()

	Convey("hi", t, func() {
		output += "done"
	})

	So(output, ShouldEqual, "done")
}

func TestSingleScopeWithMultipleConveys(t *testing.T) {
	output := prepare()

	Convey("1", t, func() {
		output += "1"
	})

	Convey("2", t, func() {
		output += "2"
	})

	So(output, ShouldEqual, "12")
}

func TestNestedScopes(t *testing.T) {
	output := prepare()

	Convey("a", t, func() {
		output += "a "

		Convey("bb", func() {
			output += "bb "

			Convey("ccc", func() {
				output += "ccc | "
			})
		})
	})

	So(output, ShouldEqual, "a bb ccc | ")
}

func TestNestedScopesWithIsolatedExecution(t *testing.T) {
	output := prepare()

	Convey("a", t, func() {
		output += "a "

		Convey("aa", func() {
			output += "aa "

			Convey("aaa", func() {
				output += "aaa | "
			})

			Convey("aaa1", func() {
				output += "aaa1 | "
			})
		})

		Convey("ab", func() {
			output += "ab "

			Convey("abb", func() {
				output += "abb | "
			})
		})
	})

	So(output, ShouldEqual, "a bb ccc | ")
}

func TestSingleScopeWithConveyAndNestedReset(t *testing.T) {
	output := prepare()

	Convey("1", t, func() {
		output += "1"

		Reset(func() {
			output += "a"
		})
	})

	So(output, ShouldEqual, "1a")
}

func TestPanicingReset(t *testing.T) {
	output := prepare()

	Convey("1", t, func() {
		output += "1"

		Reset(func() {
			panic("nooo")
		})

		Convey("runs since the reset hasn't yet", func() {
			output += "a"
		})

		Convey("but this doesnt", func() {
			output += "nope"
		})
	})

	So(output, ShouldEqual, "1a")
}

func TestSingleScopeWithMultipleRegistrationsAndReset(t *testing.T) {
	output := prepare()

	Convey("reset after each nested convey", t, func() {
		Convey("first output", func() {
			output += "1"
		})

		Convey("second output", func() {
			output += "2"
		})

		Reset(func() {
			output += "a"
		})
	})

	So(output, ShouldEqual, "1a2a")
}

func TestSingleScopeWithMultipleRegistrationsAndMultipleResets(t *testing.T) {
	output := prepare()

	Convey("each reset is run at end of each nested convey", t, func() {
		Convey("1", func() {
			output += "1"
		})

		Convey("2", func() {
			output += "2"
		})

		Reset(func() {
			output += "a"
		})

		Reset(func() {
			output += "b"
		})
	})

	So(output, ShouldEqual, "1ab2ab")
}

func Test_Failure_AtHigherLevelScopePreventsChildScopesFromRunning(t *testing.T) {
	output := prepare()

	Convey("This step fails", t, func() {
		So(1, ShouldEqual, 2)

		Convey("this should NOT be executed", func() {
			output += "a"
		})
	})

	So(output, ShouldEqual, "")
}

func Test_Panic_AtHigherLevelScopePreventsChildScopesFromRunning(t *testing.T) {
	output := prepare()

	Convey("This step panics", t, func() {
		Convey("this happens, because the panic didn't happen yet", func() {
			output += "1"
		})

		output += "a"

		Convey("this should NOT be executed", func() {
			output += "2"
		})

		output += "b"

		panic("Hi")

		output += "nope"
	})

	So(output, ShouldEqual, "1ab")
}

func Test_Panic_InChildScopeDoes_NOT_PreventExecutionOfSiblingScopes(t *testing.T) {
	output := prepare()

	Convey("This is the parent", t, func() {
		Convey("This step panics", func() {
			panic("Hi")
			output += "1"
		})

		Convey("This sibling should execute", func() {
			output += "2"
		})
	})

	So(output, ShouldEqual, "2")
}

func Test_Failure_InChildScopeDoes_NOT_PreventExecutionOfSiblingScopes(t *testing.T) {
	output := prepare()

	Convey("This is the parent", t, func() {
		Convey("This step fails", func() {
			So(1, ShouldEqual, 2)
			output += "1"
		})

		Convey("This sibling should execute", func() {
			output += "2"
		})
	})

	So(output, ShouldEqual, "2")
}

func TestResetsAreAlwaysExecutedAfterScope_Panics(t *testing.T) {
	output := prepare()

	Convey("This is the parent", t, func() {
		Convey("This step panics", func() {
			panic("Hi")
			output += "1"
		})

		Convey("This sibling step does not panic", func() {
			output += "a"

			Reset(func() {
				output += "b"
			})
		})

		Reset(func() {
			output += "2"
		})
	})

	So(output, ShouldEqual, "2ab2")
}

func TestResetsAreAlwaysExecutedAfterScope_Failures(t *testing.T) {
	output := prepare()

	Convey("This is the parent", t, func() {
		Convey("This step fails", func() {
			So(1, ShouldEqual, 2)
			output += "1"
		})

		Convey("This sibling step does not fail", func() {
			output += "a"

			Reset(func() {
				output += "b"
			})
		})

		Reset(func() {
			output += "2"
		})
	})

	So(output, ShouldEqual, "2ab2")
}

func TestSkipTopLevel(t *testing.T) {
	output := prepare()

	SkipConvey("hi", t, func() {
		output += "This shouldn't be executed!"
	})

	So(output, ShouldEqual, "")
}

func TestSkipNestedLevel(t *testing.T) {
	output := prepare()

	Convey("hi", t, func() {
		output += "yes"

		SkipConvey("bye", func() {
			output += "no"
		})
	})

	So(output, ShouldEqual, "yes")
}

func TestSkipNestedLevelSkipsAllChildLevels(t *testing.T) {
	output := prepare()

	Convey("hi", t, func() {
		output += "yes"

		SkipConvey("bye", func() {
			output += "no"

			Convey("byebye", func() {
				output += "no-no"
			})
		})
	})

	So(output, ShouldEqual, "yes")
}

func TestIterativeConveys(t *testing.T) {
	output := prepare()

	Convey("Test", t, func() {
		for x := 0; x < 10; x++ {
			y := strconv.Itoa(x)

			Convey(y, func() {
				output += y
			})
		}
	})

	So(output, ShouldEqual, "0123456789")
}

func TestClosureVariables(t *testing.T) {
	output := prepare()

	i := 0

	Convey("A", t, func() {
		i = i + 1
		j := i

		output += "A" + strconv.Itoa(i) + " "

		Convey("B", func() {
			k := j
			j = j + 1

			output += "B" + strconv.Itoa(k) + " "

			Convey("C", func() {
				output += "C" + strconv.Itoa(k) + strconv.Itoa(j) + " "
			})

			Convey("D", func() {
				output += "D" + strconv.Itoa(k) + strconv.Itoa(j) + " "
			})
		})

		Convey("C", func() {
			output += "C" + strconv.Itoa(j) + " "
		})
	})

	output += "D" + strconv.Itoa(i) + " "

	So(output, ShouldEqual, "A1 B1 C12 A2 B2 D23 A3 C3 D3 ")
}

func TestClosureVariablesWithReset(t *testing.T) {
	output := prepare()

	i := 0

	Convey("A", t, func() {
		i = i + 1
		j := i

		output += "A" + strconv.Itoa(i) + " "

		Reset(func() {
			output += "R" + strconv.Itoa(i) + strconv.Itoa(j) + " "
		})

		Convey("B", func() {
			output += "B" + strconv.Itoa(j) + " "
		})

		Convey("C", func() {
			output += "C" + strconv.Itoa(j) + " "
		})
	})

	output += "D" + strconv.Itoa(i) + " "

	So(output, ShouldEqual, "A1 B1 R11 A2 C2 R22 D2 ")
}

func TestWrappedSimple(t *testing.T) {
	prepare()
	output := resetTestString{""}

	Convey("A", t, func() {
		func() {
			output.output += "A "

			Convey("B", func() {
				output.output += "B "

				Convey("C", func() {
					output.output += "C "
				})

			})

			Convey("D", func() {
				output.output += "D "
			})
		}()
	})

	So(output, ShouldEqual, "A B C A D ", output.output)
}

type resetTestString struct {
	output string
}

func addReset(o *resetTestString, f func()) func() {
	return func() {
		Reset(func() {
			o.output += "R "
		})

		f()
	}
}

func TestWrappedReset(t *testing.T) {
	prepare()
	output := resetTestString{""}

	Convey("A", t, addReset(&output, func() {
		output.output += "A "

		Convey("B", func() {
			output.output += "B "
		})

		Convey("C", func() {
			output.output += "C "
		})
	}))

	So(output, ShouldEqual, "A B R A C R ", output.output)
}

func TestWrappedReset2(t *testing.T) {
	prepare()
	output := resetTestString{""}

	Convey("A", t, func() {
		Reset(func() {
			output.output += "R "
		})

		func() {
			output.output += "A "

			Convey("B", func() {
				output.output += "B "

				Convey("C", func() {
					output.output += "C "
				})
			})

			Convey("D", func() {
				output.output += "D "
			})
		}()
	})

	So(output, ShouldEqual, "A B C R A D R ", output.output)
}

func TestInfiniteLoopWithTrailingFail(t *testing.T) {
	done := make(chan int)

	go func() {
		Convey("This fails", t, func() {
			Convey("and this is run", func() {
				So(true, ShouldEqual, true)
			})

			/* And this prevents the whole block to be marked as run */
			So(false, ShouldEqual, true)
		})

		done <- 1
	}()

	select {
	case <-done:
		return
	case <-time.After(1 * time.Millisecond):
		t.Fail()
	}
}

func TestOutermostResetInvokedForGrandchildren(t *testing.T) {
	output := prepare()

	Convey("A", t, func() {
		output += "A "

		Reset(func() {
			output += "rA "
		})

		Convey("B", func() {
			output += "B "

			Reset(func() {
				output += "rB "
			})

			Convey("C", func() {
				output += "C "

				Reset(func() {
					output += "rC "
				})
			})

			Convey("D", func() {
				output += "D "

				Reset(func() {
					output += "rD "
				})
			})
		})
	})

	So(output, ShouldEqual, "A B C rC rB rA A B D rD rB rA ")
}

func TestFailureOption(t *testing.T) {
	output := prepare()

	Convey("A", t, FailureHalts, func() {
		output += "A "
		So(true, ShouldEqual, true)
		output += "B "
		So(false, ShouldEqual, true)
		output += "C "
	})

	So(output, ShouldEqual, "A B ")
}

func TestFailureOption2(t *testing.T) {
	output := prepare()

	Convey("A", t, func() {
		output += "A "
		So(true, ShouldEqual, true)
		output += "B "
		So(false, ShouldEqual, true)
		output += "C "
	})

	So(output, ShouldEqual, "A B ")
}

func TestFailureOption3(t *testing.T) {
	output := prepare()

	Convey("A", t, FailureContinues, func() {
		output += "A "
		So(true, ShouldEqual, true)
		output += "B "
		So(false, ShouldEqual, true)
		output += "C "
	})

	So(output, ShouldEqual, "A B C ")
}

func TestFailureOptionInherit(t *testing.T) {
	output := prepare()

	Convey("A", t, FailureContinues, func() {
		output += "A1 "
		So(false, ShouldEqual, true)
		output += "A2 "

		Convey("B", func() {
			output += "B1 "
			So(true, ShouldEqual, true)
			output += "B2 "
			So(false, ShouldEqual, true)
			output += "B3 "
		})
	})

	So(output, ShouldEqual, "A1 A2 B1 B2 B3 ")
}

func TestFailureOptionInherit2(t *testing.T) {
	output := prepare()

	Convey("A", t, FailureHalts, func() {
		output += "A1 "
		So(false, ShouldEqual, true)
		output += "A2 "

		Convey("B", func() {
			output += "A1 "
			So(true, ShouldEqual, true)
			output += "A2 "
			So(false, ShouldEqual, true)
			output += "A3 "
		})
	})

	So(output, ShouldEqual, "A1 ")
}

func TestFailureOptionInherit3(t *testing.T) {
	output := prepare()

	Convey("A", t, FailureHalts, func() {
		output += "A1 "
		So(true, ShouldEqual, true)
		output += "A2 "

		Convey("B", func() {
			output += "B1 "
			So(true, ShouldEqual, true)
			output += "B2 "
			So(false, ShouldEqual, true)
			output += "B3 "
		})
	})

	So(output, ShouldEqual, "A1 A2 B1 B2 ")
}

func TestFailureOptionNestedOverride(t *testing.T) {
	output := prepare()

	Convey("A", t, FailureContinues, func() {
		output += "A "
		So(false, ShouldEqual, true)
		output += "B "

		Convey("C", FailureHalts, func() {
			output += "C "
			So(true, ShouldEqual, true)
			output += "D "
			So(false, ShouldEqual, true)
			output += "E "
		})
	})

	So(output, ShouldEqual, "A B C D ")
}

func TestFailureOptionNestedOverride2(t *testing.T) {
	output := prepare()

	Convey("A", t, FailureHalts, func() {
		output += "A "
		So(true, ShouldEqual, true)
		output += "B "

		Convey("C", FailureContinues, func() {
			output += "C "
			So(true, ShouldEqual, true)
			output += "D "
			So(false, ShouldEqual, true)
			output += "E "
		})
	})

	So(output, ShouldEqual, "A B C D E ")
}

func TestMultipleInvocationInheritance(t *testing.T) {
	output := prepare()

	Convey("A", t, FailureHalts, func() {
		output += "A1 "
		So(true, ShouldEqual, true)
		output += "A2 "

		Convey("B", FailureContinues, func() {
			output += "B1 "
			So(true, ShouldEqual, true)
			output += "B2 "
			So(false, ShouldEqual, true)
			output += "B3 "
		})

		Convey("C", func() {
			output += "C1 "
			So(true, ShouldEqual, true)
			output += "C2 "
			So(false, ShouldEqual, true)
			output += "C3 "
		})
	})

	So(output, ShouldEqual, "A1 A2 B1 B2 B3 A1 A2 C1 C2 ")
}

func TestMultipleInvocationInheritance2(t *testing.T) {
	output := prepare()

	Convey("A", t, FailureContinues, func() {
		output += "A1 "
		So(true, ShouldEqual, true)
		output += "A2 "
		So(false, ShouldEqual, true)
		output += "A3 "

		Convey("B", FailureHalts, func() {
			output += "B1 "
			So(true, ShouldEqual, true)
			output += "B2 "
			So(false, ShouldEqual, true)
			output += "B3 "
		})

		Convey("C", func() {
			output += "C1 "
			So(true, ShouldEqual, true)
			output += "C2 "
			So(false, ShouldEqual, true)
			output += "C3 "
		})
	})

	So(output, ShouldEqual, "A1 A2 A3 B1 B2 A1 A2 A3 C1 C2 C3 ")
}

func TestSetDefaultFailureMode(t *testing.T) {
	output := prepare()

	SetDefaultFailureMode(FailureContinues) // the default is normally FailureHalts
	defer SetDefaultFailureMode(FailureHalts)

	Convey("A", t, func() {
		output += "A1 "
		So(true, ShouldBeFalse)
		output += "A2 "
	})

	So(output, ShouldEqual, "A1 A2 ")
}

func prepare() string {
	//testReporter = newNilReporter()
	return ""
}