/*
Package tmpenv is a library to overwrite environment variables temporarily and to restore them
with backup. All modified environment variable is restored when Restore() method is called.
And all added environment variables are removed when Restore() method is called.

Setting temporary environment variables and restoring them are common pattern on testing.
However, they have a pitfall since setting empty value to environment variable is different from
unsetting environment variable.

tmpenv helps you to do the pattern correctly with less lines.

This library is usually used for testing with environment variables.
*/
package tmpenv
