# Simple scrapper

Includes concurrency and simple cache.
Made as an interview assignment.

Output:
```bash
requesting https://pl.wikipedia.org/wiki/Niemcy
requesting https://pl.wikipedia.org/wiki/Polska
requesting https://pl.wikipedia.org/wiki/Francja
done https://pl.wikipedia.org/wiki/Niemcy [{na 331} {i 459} {w 920}]
requesting https://pl.wikipedia.org/wiki/Polska
done https://pl.wikipedia.org/wiki/Polska [{z 442} {i 786} {w 1324}]
from cache https://pl.wikipedia.org/wiki/Polska
from cache https://pl.wikipedia.org/wiki/Polska
from cache https://pl.wikipedia.org/wiki/Polska
done https://pl.wikipedia.org/wiki/Francja [{i 159} {ΓÇô 203} {w 273}]
done https://pl.wikipedia.org/wiki/Polska [{z 442} {i 786} {w 1324}]
done https://pl.wikipedia.org/wiki/Polska [{z 442} {i 786} {w 1324}]
done https://pl.wikipedia.org/wiki/Polska [{z 442} {i 786} {w 1324}]
done https://pl.wikipedia.org/wiki/Polska [{z 442} {i 786} {w 1324}]
```