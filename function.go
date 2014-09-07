
package interpreter

import "github.com/Cirru/parser"

type context struct {
  env *Env
  args []interface{}
  code []interface{}
}

func (env *Env) fn(xs []interface{}) (ret unitype) {
  ret.Type = cirruFn
  if args, ok := xs[0].([]interface{}); ok {
    ret.Value = context{env, args, xs[1:]}
  }
  return
}

func (env *Env) call(xs []interface{}) (ret unitype) {
  fn := env.get(xs[0:1])
  if fn.Type == cirruFn {
    if item, ok := fn.Value.(context); ok {
      runtime := Env{}
      for i, para := range item.args {
        // println("i is:", i)
        if token, ok := para.(parser.Token); ok {
          runtime[token.Text] = env.get(xs[i+1:i+2])
        }
      }
      for _, line := range item.code {
        if exp, ok := line.([]interface{}); ok {
          ret = Evaluate(&runtime, exp)
        }
      }
      return
    }
  } else {
    panic("not function")
  }
  return
}