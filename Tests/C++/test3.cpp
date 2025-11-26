#define BOOST_TEST_MODULE MyTests

#include <boost/test/included/unit_test.hpp>
#include "Array.h"
#include "DoubleList.h"
#include "ForwardList.h"
#include "Queue.h"
#include "Stack.h"
#include "ChainMap.h"
#include "BlackRedTree.h"
#include <stdexcept>

BOOST_AUTO_TEST_SUITE(ArrayTests)

BOOST_AUTO_TEST_CASE(ArrayConstructorAndBasicOperations) {
    Array arr(5);
    BOOST_CHECK_EQUAL(arr.getLength(), 0);
    BOOST_CHECK_EQUAL(arr.getSize(), 5);
    
    arr.addElementEnd("first");
    arr.addElementEnd("second");
    BOOST_CHECK_EQUAL(arr.getLength(), 2);
    BOOST_CHECK_EQUAL(arr.getElement(0), "first");
    BOOST_CHECK_EQUAL(arr.getElement(1), "second");
}

BOOST_AUTO_TEST_CASE(ArrayStringConstructor) {
    Array arr("test");
    BOOST_CHECK_EQUAL(arr.getLength(), 4);
    BOOST_CHECK_EQUAL(arr.getElement(0), "t");
    BOOST_CHECK_EQUAL(arr.getElement(3), "t");
}

BOOST_AUTO_TEST_CASE(ArrayAddAndDelete) {
    Array arr(3);
    arr.addElementEnd("one");
    arr.addElementEnd("two");
    arr.addElementEnd("three");
    
    BOOST_CHECK_EQUAL(arr.getLength(), 3);
    
    arr.deleteElement(1);
    BOOST_CHECK_EQUAL(arr.getLength(), 2);
    BOOST_CHECK_EQUAL(arr.getElement(0), "one");
    BOOST_CHECK_EQUAL(arr.getElement(1), "three");
    
    arr.addElementAtIndex("new", 1);
    BOOST_CHECK_EQUAL(arr.getElement(1), "new");
    BOOST_CHECK_EQUAL(arr.getLength(), 3);
}

BOOST_AUTO_TEST_CASE(ArraySearch) {
    Array arr(5);
    arr.addElementEnd("apple");
    arr.addElementEnd("banana");
    arr.addElementEnd("cherry");
    
    BOOST_CHECK_EQUAL(arr.isInArray("banana"), 1);
    BOOST_CHECK_EQUAL(arr.isInArray("grape"), -1);
}

BOOST_AUTO_TEST_CASE(ArrayExceptions) {
    Array arr(2);
    BOOST_CHECK_THROW(Array(0), runtime_error);
    BOOST_CHECK_THROW(arr.getElement(0), range_error);
    BOOST_CHECK_THROW(arr.setElement("test", 5), range_error);
    BOOST_CHECK_THROW(arr.deleteElement(0), range_error);
    BOOST_CHECK_THROW(arr.addElementAtIndex("test", 5), range_error);
}

BOOST_AUTO_TEST_CASE(ArraySetElement) {
    Array arr(3);
    arr.addElementEnd("one");
    arr.addElementEnd("two");
    arr.setElement("modified", 1);
    BOOST_CHECK_EQUAL(arr.getElement(1), "modified");
}

BOOST_AUTO_TEST_CASE(ArrayExtend) {
    Array arr(2);
    arr.addElementEnd("one");
    arr.addElementEnd("two");
    arr.addElementEnd("three");
    BOOST_CHECK(arr.getSize() > 2);
    BOOST_CHECK_EQUAL(arr.getElement(2), "three");
}

BOOST_AUTO_TEST_CASE(ArrayGetHead) {
    Array arr(3);
    arr.addElementEnd("test");
    ArNode* head = arr.getHead();
    BOOST_CHECK(head != nullptr);
}

BOOST_AUTO_TEST_SUITE_END()

BOOST_AUTO_TEST_SUITE(DoubleListTests)

BOOST_AUTO_TEST_CASE(DoubleListAddAndRemove) {
    DoubleList list;
    list.addHead("first");
    list.addTail("third");
    list.addAfter("second", 0);
    
    BOOST_CHECK_EQUAL(list.getLength(), 3);
    BOOST_CHECK_EQUAL(list.getElement(0), "first");
    BOOST_CHECK_EQUAL(list.getElement(1), "second");
    BOOST_CHECK_EQUAL(list.getElement(2), "third");
    
    list.deleteAt(1);
    BOOST_CHECK_EQUAL(list.getLength(), 2);
    BOOST_CHECK_EQUAL(list.getElement(0), "first");
    BOOST_CHECK_EQUAL(list.getElement(1), "third");
}

BOOST_AUTO_TEST_CASE(DoubleListFindAndPop) {
    DoubleList list;
    list.addHead("apple");
    list.addTail("banana");
    list.addTail("cherry");
    
    DFNode* found = list.findByValue("banana");
    BOOST_CHECK(found != nullptr);
    
    string popped = list.popElement(1);
    BOOST_CHECK_EQUAL(popped, "banana");
    BOOST_CHECK_EQUAL(list.getLength(), 2);
    
    BOOST_CHECK(list.findByValue("nonexistent") == nullptr);
}

BOOST_AUTO_TEST_CASE(DoubleListEdgeCases) {
    DoubleList list;
    BOOST_CHECK(list.isEmpty());
    
    list.addHead("only");
    BOOST_CHECK_EQUAL(list.getLength(), 1);
    
    list.deleteHead();
    BOOST_CHECK(list.isEmpty());
    
    list.addTail("new");
    list.deleteTail();
    BOOST_CHECK(list.isEmpty());
}

BOOST_AUTO_TEST_CASE(DoubleListAddBeforeAfter) {
    DoubleList list;
    list.addHead("middle");
    list.addBefore("before", 0);
    list.addAfter("after", 1);
    
    BOOST_CHECK_EQUAL(list.getLength(), 3);
    BOOST_CHECK_EQUAL(list.getElement(0), "before");
    BOOST_CHECK_EQUAL(list.getElement(1), "middle");
    BOOST_CHECK_EQUAL(list.getElement(2), "after");
}

BOOST_AUTO_TEST_CASE(DoubleListDeleteByValue) {
    DoubleList list;
    list.addTail("one");
    list.addTail("two");
    list.addTail("three");
    
    list.deleteByValue("two");
    BOOST_CHECK_EQUAL(list.getLength(), 2);
    BOOST_CHECK_EQUAL(list.getElement(0), "one");
    BOOST_CHECK_EQUAL(list.getElement(1), "three");
    
    BOOST_CHECK_THROW(list.deleteByValue("nonexistent"), runtime_error);
}

BOOST_AUTO_TEST_CASE(DoubleListExceptions) {
    DoubleList list;
    BOOST_CHECK_THROW(list.getElement(0), out_of_range);
    BOOST_CHECK_THROW(list.popElement(0), out_of_range);
    BOOST_CHECK_THROW(list.deleteAt(0), out_of_range);
}

BOOST_AUTO_TEST_CASE(DoubleListFromTailTraversal) {
    DoubleList list;
    for(int i = 0; i < 10; i++) {
        list.addTail("node" + to_string(i));
    }
    BOOST_CHECK_EQUAL(list.getElement(9), "node9");
    BOOST_CHECK_EQUAL(list.getElement(5), "node5");
}

BOOST_AUTO_TEST_SUITE_END()

BOOST_AUTO_TEST_SUITE(ForwardListTests)

BOOST_AUTO_TEST_CASE(ForwardListBasicOperations) {
    ForwardList list;
    list.pushFront("first");
    list.pushBack("third");
    list.insertBefore("second", 1);
    
    BOOST_CHECK_EQUAL(list.getSize(), 3);
    BOOST_CHECK_EQUAL(list.front(), "first");
    BOOST_CHECK_EQUAL(list.back(), "third");
    BOOST_CHECK_EQUAL(list.getAt(1), "second");
}

BOOST_AUTO_TEST_CASE(ForwardListInsertRemove) {
    ForwardList list;
    list.pushBack("one");
    list.pushBack("two");
    list.pushBack("three");
    
    list.insertAfter("new", 1);
    BOOST_CHECK_EQUAL(list.getAt(2), "new");

    list.insertBefore("head", 0);
    BOOST_CHECK_EQUAL(list.getAt(0), "head");
    
    list.removeByValue("two");
    BOOST_CHECK_EQUAL(list.getSize(), 4);
    BOOST_CHECK_EQUAL(list.getAt(1), "one");
    
    list.popFront();
    BOOST_CHECK_EQUAL(list.front(), "one");
    
    list.popBack();
    BOOST_CHECK_EQUAL(list.back(), "new");
}

BOOST_AUTO_TEST_CASE(ForwardListFind) {
    ForwardList list;
    list.pushBack("apple");
    list.pushBack("banana");
    list.pushBack("cherry");
    
    FWNode* found = list.findByValue("banana");
    BOOST_CHECK(found != nullptr);
    BOOST_CHECK_EQUAL(found->key, "banana");
    
    found = list.findByValue("nonexistent");
    BOOST_CHECK(found == nullptr);
}

BOOST_AUTO_TEST_CASE(ForwardListRemoveAfter) {
    ForwardList list;
    list.pushBack("one");
    list.pushBack("two");
    list.pushBack("three");
    
    FWNode* first = list.findByValue("one");
    list.removeAfter(first);
    
    BOOST_CHECK_EQUAL(list.getSize(), 2);
    BOOST_CHECK_EQUAL(list.getAt(0), "one");
    BOOST_CHECK_EQUAL(list.getAt(1), "three");
}

BOOST_AUTO_TEST_CASE(ForwardListExceptions) {
    ForwardList list;
    BOOST_CHECK_THROW(list.front(), runtime_error);
    BOOST_CHECK_THROW(list.back(), runtime_error);
    BOOST_CHECK_THROW(list.getAt(0), out_of_range);
    BOOST_CHECK_THROW(list.popFront(), runtime_error);
    BOOST_CHECK_THROW(list.popBack(), runtime_error);
    
    list.pushBack("test");
    FWNode* node = nullptr;
    BOOST_CHECK_THROW(list.removeAfter(node), invalid_argument);
}

BOOST_AUTO_TEST_CASE(ForwardListEmptyOperations) {
    ForwardList list;
    BOOST_CHECK(list.isEmpty());
    BOOST_CHECK_EQUAL(list.getSize(), 0);
    BOOST_CHECK(!list.removeByValue("nonexistent"));
}

BOOST_AUTO_TEST_SUITE_END()

BOOST_AUTO_TEST_SUITE(QueueTests)

BOOST_AUTO_TEST_CASE(QueueEnqueueDequeue) {
    Queue queue(0);
    queue.enqueue("first");
    queue.enqueue("second");
    queue.enqueue("third");
    
    BOOST_CHECK_EQUAL(queue.getSize(), 3);
    
    string item = queue.dequeue();
    BOOST_CHECK_EQUAL(item, "first");
    BOOST_CHECK_EQUAL(queue.getSize(), 2);
}

BOOST_AUTO_TEST_CASE(QueueDelete) {
    Queue queue(0);
    queue.enqueue("apple");
    queue.enqueue("banana");
    queue.enqueue("cherry");
    
    queue.del("banana");
    BOOST_CHECK_EQUAL(queue.getSize(), 2);
    
    queue.del("apple");
    BOOST_CHECK_EQUAL(queue.getHead()->data, "cherry");

    queue.del("cherry");
    BOOST_CHECK_EQUAL(queue.getSize(), 0);
}

BOOST_AUTO_TEST_CASE(QueueEdgeCases) {
    Queue queue(0);
    BOOST_CHECK_EQUAL(queue.getSize(), 0);
    
    queue.enqueue("single");
    BOOST_CHECK_EQUAL(queue.getSize(), 1);
    BOOST_CHECK_EQUAL(queue.dequeue(), "single");
    BOOST_CHECK_EQUAL(queue.getSize(), 0);
    
    queue.enqueue("test");
    queue.del("nonexistent");
    BOOST_CHECK_EQUAL(queue.getSize(), 1);
}

BOOST_AUTO_TEST_CASE(QueueExceptions) {
    Queue queue(0);
    BOOST_CHECK_THROW(queue.dequeue(), underflow_error);
    
    for(size_t i = 0; i < 1000; i++) {
        queue.enqueue("item" + to_string(i));
    }
    BOOST_CHECK_THROW(queue.enqueue("overflow"), overflow_error);
}

BOOST_AUTO_TEST_CASE(QueueGetHead) {
    Queue queue(0);
    queue.enqueue("test");
    QNode* head = queue.getHead();
    BOOST_CHECK(head != nullptr);
}

BOOST_AUTO_TEST_SUITE_END()

BOOST_AUTO_TEST_SUITE(StackTests)

BOOST_AUTO_TEST_CASE(StackPushPop) {
    Stack stack;
    stack.push(1);
    stack.push(2);
    
    BOOST_CHECK_EQUAL(stack.getSize(), 2);
    BOOST_CHECK_EQUAL(stack.pop(), 2);
    BOOST_CHECK_EQUAL(stack.pop(), 1);
    BOOST_CHECK(stack.isEmpty());
}

BOOST_AUTO_TEST_CASE(StackOverflowUnderflow) {
    Stack stack;
    stack.push(1);
    stack.push(2);
    
    BOOST_CHECK_THROW(stack.push(3), overflow_error);
    
    stack.pop();
    stack.pop();
    BOOST_CHECK_THROW(stack.pop(), underflow_error);
}

BOOST_AUTO_TEST_CASE(StackEmptyOperations) {
    Stack stack;
    BOOST_CHECK(stack.isEmpty());
    BOOST_CHECK_EQUAL(stack.getSize(), 0);
}

BOOST_AUTO_TEST_SUITE_END()

BOOST_AUTO_TEST_SUITE(ChainMapTests)

BOOST_AUTO_TEST_CASE(ChainMapAddFind) {
    ChainMap map(10);
    map.add("key1", 100);
    map.add("key2", 200);
    
    BOOST_CHECK(map.isContain("key1"));
    BOOST_CHECK_EQUAL(map.find("key1"), 100);
    BOOST_CHECK_EQUAL(map.find("key2"), 200);
}

BOOST_AUTO_TEST_CASE(ChainMapDelete) {
    ChainMap map(10);
    map.add("apple", 1);
    map.add("banana", 2);
    
    map.del("apple");
    BOOST_CHECK(!map.isContain("apple"));
    BOOST_CHECK(map.isContain("banana"));
    BOOST_CHECK_EQUAL(map.find("banana"), 2);
}

BOOST_AUTO_TEST_CASE(ChainMapRehash) {
    ChainMap map(2);
    map.add("key1", 1);
    map.add("key2", 2);
    map.add("key3", 3);
    
    BOOST_CHECK(map.isContain("key1"));
    BOOST_CHECK(map.isContain("key2"));
    BOOST_CHECK(map.isContain("key3"));
    BOOST_CHECK_EQUAL(map.find("key3"), 3);
}

BOOST_AUTO_TEST_CASE(ChainMapUpdateValue) {
    ChainMap map(10);
    map.add("key", 100);
    map.add("key", 200);
    
    BOOST_CHECK_EQUAL(map.find("key"), 200);
}

BOOST_AUTO_TEST_CASE(ChainMapGetAllKeys) {
    ChainMap map(10);
    map.add("a", 1);
    map.add("b", 2);
    map.add("c", 3);
    
    ChainMap keys(10);
    map.getAllKeys(keys);
    
    BOOST_CHECK(keys.isContain("a"));
    BOOST_CHECK(keys.isContain("b"));
    BOOST_CHECK(keys.isContain("c"));
    
    string allKeys = map.getAllKeysAsString();
    BOOST_CHECK(!allKeys.empty());
}

BOOST_AUTO_TEST_CASE(ChainMapExceptions) {
    ChainMap map(10);
    BOOST_CHECK_THROW(map.find("nonexistent"), runtime_error);
}

BOOST_AUTO_TEST_CASE(ChainMapPrintContents) {
    ChainMap map(5);
    map.add("test1", 1);
    map.add("test2", 2);
    map.printContents();
}

BOOST_AUTO_TEST_SUITE_END()

BOOST_AUTO_TEST_SUITE(RBTreeTests)

BOOST_AUTO_TEST_CASE(firstFixViolationTest) { //Left parent, red uncle
    RBTree st;
    st.insert(13);
    st.insert(5);
    st.insert(26);
    st.insert(4);
}

BOOST_AUTO_TEST_CASE(secondFixViolationTest) { //Left parent, black uncle, left child
    RBTree st;
    st.insert(13);
    st.insert(5);
    st.insert(26);
    st.insert(4);
    st.insert(3);
}

BOOST_AUTO_TEST_CASE(thirdFixViolationTest) { //Left parent, black uncle, right child
    RBTree st;
    st.insert(13);
    st.insert(7);
    st.insert(26);
    st.insert(4);
    st.insert(6);
}

BOOST_AUTO_TEST_CASE(sixthFixViolationTest) { //Right parent, black uncle, right child
    RBTree st;
    st.insert(13);
    st.insert(5);
    st.insert(26);
    st.insert(30);
    st.insert(28);
}

BOOST_AUTO_TEST_CASE(RBTreeComplexOperations) {
    RBTree tree;
    for(int i = 0; i < 10; i++) {
        tree.insert(i);
    }
    
    for(int i = 0; i < 10; i += 2) {
        tree.del(i);
    }
    
    for(int i = 1; i < 10; i += 2) {
        BOOST_CHECK_EQUAL(tree.get(i), i);
    }
    
    for(int i = 0; i < 10; i += 2) {
        BOOST_CHECK_THROW(tree.get(i), runtime_error);
    }
}

BOOST_AUTO_TEST_CASE(RBTreeEmptyOperations) {
    RBTree tree;
    BOOST_CHECK_THROW(tree.get(1), runtime_error);
    BOOST_CHECK_THROW(tree.del(1), runtime_error);
}

BOOST_AUTO_TEST_CASE(RBTreeDuplicateInsert) {
    RBTree tree;
    tree.insert(10);
    BOOST_CHECK_THROW(tree.insert(10), runtime_error);
}

BOOST_AUTO_TEST_CASE(RBTreeRootDeletion) {
    RBTree tree;    
    tree.insert(5);
    tree.insert(3);
    tree.insert(7);
    tree.del(5);
    BOOST_CHECK_EQUAL(tree.get(3), 3);
    BOOST_CHECK_EQUAL(tree.get(7), 7);
}

BOOST_AUTO_TEST_CASE(RBTreeComplexDeletionScenarios) {
    RBTree tree;
    tree.insert(50);
    tree.insert(30);
    tree.insert(70);
    tree.insert(20);
    tree.insert(40);
    tree.insert(60);
    tree.insert(80);
    
    tree.del(30);
    BOOST_CHECK_THROW(tree.get(30), runtime_error);
    BOOST_CHECK_THROW(tree.del(30), runtime_error);
    BOOST_CHECK_EQUAL(tree.get(20), 20);
    
    tree.del(40);
    BOOST_CHECK_THROW(tree.get(40), runtime_error);

    tree.insert(55);
    tree.del(20);
    BOOST_CHECK_THROW(tree.get(20), runtime_error);

    tree.del(80);
    BOOST_CHECK_THROW(tree.get(80), runtime_error);
}

BOOST_AUTO_TEST_SUITE_END()