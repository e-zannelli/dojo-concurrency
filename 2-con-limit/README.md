La fonction process doit appeler la fonction Call du client pour chaque élément de generator.
Les réponses doivent être renvoyées en retour de la fonction.

Le client envoyé en paramètre a une limite de 10 appels simultanés, les appels en trop seront ralentis.

Modifier le code pour faire les appels de manière concurrente, en respectant la limite de 10 appels simultanés.