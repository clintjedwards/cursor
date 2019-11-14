# Cursor

Cursor is a CI/CD/Cron-like system that uses hashicorp's go-plugin system to allow users to write their "thing do-er" pipelines in an actual programming language like golang and not some hacky configuration language

It's a large undertaking and i'll probably never finish it though

Pipeline runs cannot be individually cleared they can be all cleared or none at all. Its a history. This is because its a simplier datastructure for ordering to disallow individual deletes.
