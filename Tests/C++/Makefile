all: clean lcov launch comp
	genhtml coverage.info --ignore-errors inconsistent --ignore-errors category --output-directory coverage_html

lcov: launch
	lcov --ignore-errors inconsistent --capture --directory . --output-file coverage.info

launch: comp
	./Test

comp:
	g++ -g -O0 --coverage -o Test *.cpp

clean:
	rm -rf Test* coverage*