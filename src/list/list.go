package list

import (

)

type List interface {
    Len() int
    At(i int) interface{}
    
    //PushFront(elm interface{})
    //PushBack(elm interface{})
    
    //RemoveAt(i int)
}