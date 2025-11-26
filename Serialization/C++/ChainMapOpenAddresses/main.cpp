#include "OpenAddressMap.h"

int main() {
    OpenAddressingMap a(3);
    a.add("apples", 3);
    a.add("bubbles", 5);
    a.add("gems", 2);
    a.printContents();
    a.writeBinary("test");
    OpenAddressingMap b(1);
    b.readBinary("test");
    b.printContents();

    OpenAddressingMap c(3);
    c.add("apples", 3);
    c.add("bubbles", 5);
    c.add("gems", 2);
    c.printContents();
    c.writeText("test2");
    OpenAddressingMap d(1);
    d.readText("test2");
    d.printContents();

    return 0;
}