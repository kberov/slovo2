# Slovo2 -- Искони бѣ Слово (At the beginning was the Word)

I still remember the time somewhwere back at year 2000 when I was waiting a
whole evening to download ActiveState Perl 5.6.1 (or it was 5.6.0?…) with the
intention to write a [CGI
script](https://en.wikipedia.org/wiki/Common_Gateway_Interface). I've chosen
Perl back then, because it was hard (they said) to write CGI scripts in C. PHP
and Python were not yet considered mature enough in their evolution.

After learning some and definitely after reading [Beginning
Perl](https://www.perl.org/books/beginning-perl/) in Bulgarian, I was able to
more independently take care of my growing family and untill very recently I
did and mangled for living mostly Perl programs. It is still possible today.
Some huge company still use Perl and make a lot of money. And it is still in top
10 languages as of this writing 2025-07-27 at
[TIOBE](https://www.tiobe.com/tiobe-index/) … ook, it's eleventh now.

Ok, but some years aGo I stumbled on something even easier to learn from
scratch asa beginner and tens of times faster -- [Go](https://go.dev/).

Slovo2 is reimplementation (this time in Go)
[Slovo](https://github.com/kberov/Slovo) -- my previous lonely project, which
by the way runs our small family business -- a publishing house book store:
https://слово.бг . The database is completely the same -- one single
[SQLite](https://sqlite.org/whentouse.html) file.

The application is very similar in spirit of the
[Mojolicious](https://mojolicious.org/) application and framework Slovo. Only
the syntax is simpler and between 20 and 30 times faster even in CGI mode. Yes
CGI is viable today and even the cheapest shared hosting meets all
requirements of Slovo2.

## Features

It is still alpha quality... TODO _ write a list of features...
