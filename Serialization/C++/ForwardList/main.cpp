#include "ForwardList.h"

int main() {
    ForwardList a = {"HELP", "I CANT", "HOLD IT", "ANYMORE"};
    a.print();
    a.writeBinary("test");
    ForwardList b;
    b.readBinary("test");
    b.print();

    ForwardList c = {"HELP", "I CANT", "HOLD IT", "ANYMORE"};
    c.print();
    c.writeText("test2");
    ForwardList d;
    d.readText("test2");
    d.print();

    return 0;
}