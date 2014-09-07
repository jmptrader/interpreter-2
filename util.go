
package interpreter

import (
  "github.com/Cirru/parser"
  "github.com/Cirru/writer"
  "fmt"
)

func stringifyunitype(data unitype) string {
  tree := transformunitype(data)
  lines := []interface{}{tree}
  return writer.MakeCode(lines)
}

func transformunitype(data unitype) []interface{} {
  switch data.Type {
    case cirruString:
      if stringValue, ok := data.Value.(string); ok {
        return []interface{}{"string", stringValue}
      }
    case cirruInt:
      if value, ok := data.Value.(int); ok {
        str := fmt.Sprintf("%d", value)
        return []interface{}{"int", str}
      }
    case cirruFloat:
      if value, ok := data.Value.(float64); ok {
        str := fmt.Sprintf("%g", value)
        return []interface{}{"float", str}
      }
    case cirruBool:
      if value, ok := data.Value.(bool); ok {
        if value {
          return []interface{}{"bool", "true"}
        }
        return []interface{}{"bool", "false"}
      }
    case cirruArray:
      list := []interface{}{"array"}
      if value, ok := data.Value.(*[]unitype); ok {
        for _, item := range *value {
          list = append(list, transformunitype(item))
        }
      }
      return list
    case cirruTable:
      list := []interface{}{"table"}
      if value, ok := data.Value.(*Env); ok {
        for k, v := range *value {
          pair := []interface{}{k, transformunitype(v)}
          list = append(list, pair)
        }
      }
      return list
    case cirruRegexp:
      str := fmt.Sprintf("%s", data.Value)
      return []interface{}{"regexp", str}
    case cirruFn:
      if fnContext, ok := data.Value.(context); ok {
        args := transformCode(fnContext.args)
        code := transformCode(fnContext.code)
        return []interface{}{"fn", args, code}
      }
    case cirruNil:
      return []interface{}{"nil"}
    default:
      panic("unknown structure")
  }
  return []interface{}{}
}

func generateString(x string) (ret unitype) {
  ret.Type = cirruString
  ret.Value = x
  return
}

func generateTable(m *Env) (ret unitype) {
  ret.Type = cirruTable
  ret.Value = m
  return
}

func transformCode(xs []interface{}) []interface{} {
  hold := []interface{}{}
  for _, item := range xs {
    if buffer, ok := item.(parser.Token); ok {
      hold = append(hold, buffer.Text)
    }
    if list, ok := item.([]interface{}); ok {
      hold = append(hold, transformCode(list))
    }
  }
  return hold
}
