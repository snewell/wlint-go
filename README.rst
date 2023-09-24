wlint-go
========
|codacy|
|code-climate|
|go-report-card|

This is a rewrite of `wlint`_ in Go.  The goal is to create a series of linters
to help with writing projects.

Installing
----------
:code:`wlint` is only released via source.  Make sure you have an
up-to-date Go installed and run the following in your terminal:

.. code-block:: bash

    $ go install github.com/snewell/wlint-go
    $

If that works, you'll have :code:`wlint-go` as an executable in
:code:`$GOROOT/bin`.  If you prefer the name :code:`wlint`, the included
:code:`Makefile` will do the trick.

.. code-block:: bash

    /path/to/wlint/src $ make
    /path/to/wlint/src $ make install

You can specify an optional destination directory by setting
:code:`DESTDIR` during the call to :code:`make install`.

.. code-block:: bash

    /path/to/wlint/src $ DESTDIR=/usr/local/bin make install

.. _wlint: https://github.com/snewell/wlint

.. |codacy| image:: https://app.codacy.com/project/badge/Grade/153bdcd317c04cd1aefcaa937eb35011
    :target: https://app.codacy.com/gh/snewell/wlint-go/dashboard?utm_source=gh&utm_medium=referral&utm_content=&utm_campaign=Badge_grade
.. |code-climate| image:: https://api.codeclimate.com/v1/badges/01bd605f255abaf54b12/maintainability
   :target: https://codeclimate.com/github/snewell/wlint-go/maintainability
   :alt: Maintainability
.. |go-report-card| image:: https://goreportcard.com/badge/github.com/snewell/wlint-go
    :target: https://goreportcard.com/report/github.com/snewell/wlint-go
    :alt: Go Report Card
