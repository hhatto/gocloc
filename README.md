# gocloc

A little fast cloc(Count Lines Of Code), written in Go.
Inspired by [tokei](https://github.com/Aaronepower/tokei).

**This is experimental module. Highly under development.**

## Installation

```
$ go get github.com/hhatto/gocloc
```

## Usage

```
$ gocloc .
```

```
$ gocloc .
-------------------------------------------------------------------------------
Language                     files          blank        comment           code
-------------------------------------------------------------------------------
Markdown                         3              8              0             18
Go                               1             29              1            323
-------------------------------------------------------------------------------
TOTAL                            4             37              1            341
-------------------------------------------------------------------------------
```

## Performance
* cloc 1.66
* tokei 1.5.1
* target repository is [golang/go commit:633ab74](https://github.com/golang/go/tree/633ab7426a906b72dcf6f1d54e87f4ae926dc4e1)

### cloc

```
$ time cloc .
    5171 text files.
    5052 unique files.
     420 files ignored.

https://github.com/AlDanial/cloc v 1.66  T=23.31 s (204.0 files/s, 48203.3 lines/s)
-----------------------------------------------------------------------------------
Language                         files          blank        comment           code
-----------------------------------------------------------------------------------
Go                                4197         101140         125939         780280
Assembly                           331           6128          14654          40531
HTML                                41           4927            198          33316
C                                   90           1076            940           7390
Perl                                12            185            177           1135
Bourne Again Shell                  25            209            266            933
XML                                  4             85              9            623
Bourne Shell                         8             71            302            467
Python                               1            121             88            295
DOS Batch                            5             55              1            238
JavaScript                           4             48            122            231
C/C++ Header                        15             50            147            211
CSS                                  3             51              9            176
yacc                                 1             27             20            155
Protocol Buffers                     1              1              0            144
Windows Resource File                4             25              0            116
JSON                                 2              0              0             36
make                                 7             12             10             35
Fortran 90                           2              1              3              8
C++                                  1              3              5              7
awk                                  1              1              6              7
-----------------------------------------------------------------------------------
SUM:                              4755         114216         142896         866334
-----------------------------------------------------------------------------------
cloc .  14.98s user 9.47s system 103% cpu 23.697 total
```

### tokei

```
$ tokei .
-------------------------------------------------------------------------------
 Language            Files        Total       Blanks     Comments         Code
-------------------------------------------------------------------------------
 BASH                   27         2134          260          570         1304
 Batch                   5          294           55            0          239
 C                      92         9436         1081          946         7409
 C++                     1           15            3            5            7
 CSS                     3          236           51            9          176
 FORTRAN Modern          2           12            1            3            8
 Go                   4272      1027537       103241       150411       773970
 C Header               15          408           50          147          211
 HTML                   41        38441         4927          204        33316
 JavaScript              4          401           48          122          231
 JSON                    2            0            0            0           36
 Makefile                7           57           13           10           34
 Markdown                3            0            0            0          115
 Perl                   10         1255          151         1096          343
 Protocol Buffers        1          145            1            0          144
 Python                  1          504          121           56          327
 Assembly              334        61318         6130            0        55188
 Plain Text             28            0            0            0       137715
 XML                     4          717           85            9          623
-------------------------------------------------------------------------------
 Total                4852      1142910       116218       153588      1011396
-------------------------------------------------------------------------------
tokei .  1.30s user 0.06s system 99% cpu 1.358 total
```

### gocloc

```
$ gocloc .
-------------------------------------------------------------------------------
Language                     files          blank        comment           code
-------------------------------------------------------------------------------
FORTRAN Modern                   2              1              3              8
JavaScript                       4             48            122            231
JSON                             2              0              0             36
Awk                              1              1              6              7
C++                              1              3              5              7
Batch                            5             55              0            239
CSS                              3             51              9            176
Yacc                             1             27             20            155
Perl                            10            151            158            946
Assembly                       329           6092              0          54906
C                               92           1081            996           7409
XML                              4             85              9            623
Protocol Buffers                 1              1              0            144
Go                            4272         103241         136154         788685
C Header                        15             50            147            211
Python                           1            121             56            327
Markdown                         3             29              0             86
BASH                            25            254            558           1294
Plain Text                      28           2066              0         135647
HTML                            41           4927            198          33316
-------------------------------------------------------------------------------
TOTAL                         4840         118284         138441        1024453
-------------------------------------------------------------------------------
gocloc .  0.49s user 0.06s system 105% cpu 0.524 total
```

## License
MIT
