# Simple scrapper

Includes concurrency and simple cache.
Made as an interview assignment.

Output:
```bash
requesting https://pl.wikipedia.org/wiki/Niemcy
requesting https://pl.wikipedia.org/wiki/Polska
requesting https://pl.wikipedia.org/wiki/Francja
done https://pl.wikipedia.org/wiki/Polska [{na 381} {z 442} {i 786} {w 1324}]
from cache https://pl.wikipedia.org/wiki/Polska
from cache https://pl.wikipedia.org/wiki/Polska
from cache https://pl.wikipedia.org/wiki/Polska
from cache https://pl.wikipedia.org/wiki/Polska
done https://pl.wikipedia.org/wiki/Niemcy [{z 315} {na 331} {i 459} {w 920}]
done https://pl.wikipedia.org/wiki/Polska [{na 381} {z 442} {i 786} {w 1324}]
done https://pl.wikipedia.org/wiki/Polska [{na 381} {z 442} {i 786} {w 1324}]
done https://pl.wikipedia.org/wiki/Polska [{na 381} {z 442} {i 786} {w 1324}]
done https://pl.wikipedia.org/wiki/Polska [{na 381} {z 442} {i 786} {w 1324}]
done https://pl.wikipedia.org/wiki/Francja [{edytuj 107} {na 110} {i 159} {w 273}]
```