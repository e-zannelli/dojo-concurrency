La fonction process doit appeler la fonction Call du client pour chaque élément de generator.
Les réponses doivent être renvoyées en retour de la fonction.

Un cache mémoire doit être utilisé pour stocker les réponses, si une réponse est déjà présente dans le cache, elle doit être renvoyée sans appeler le client.

Les accès concurrents au cache doivent être gérés correctement.
