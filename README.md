# Comparaison Haute Performance : WebSocket vs gRPC

## C'est quoi ce projet ?

J'ai créé une simulation simple : des Pokémon qui bougent et interagissent dans une arène virtuelle.

L'idée est de voir comment WebSocket et gRPC gèrent l'envoi rapide et continu d'informations (position, état des Pokémon), comme dans un jeu ou une application live.

## Comment ça marche ? 

* **Clients Pokémon (Go)** : Simulent les Pokémon. Ils envoient leurs infos via WebSocket ou gRPC.
* **Serveur WebSocket (Go)** : Reçoit les données des clients WS, les centralise et les envoie au dashboard.
* **Serveur gRPC (Go)** : Reçoit les données des clients gRPC via un flux continu.Dans cet exemple, il affiche surtout les données reçues dans ses logs.
* **Dashboard (Ebiten/Go)** : Se connecte au serveur WebSocket pour afficher les Pokémon correspondants sur une interface graphique.

## Le cœur du système : Les applications

### Les Clients Pokémon (en Go)

Chaque client représente un Pokémon :
* Se déplace aléatoirement.
* Son état (PV, action) change au hasard.
* Envoie ses nouvelles informations (position, PV, action, et un timestamp d'envoi T1) à son serveur dédié (WebSocket ou gRPC) à intervalle régulier (configurable via `UPDATE_INTERVAL_MS`).
* Mesure la latence aller-retour (RTT) en calculant la différence entre le moment où il reçoit l'écho de ses propres données (T3) et le moment où il les avait envoyées (T1).

### Le Serveur WebSocket (en Go)

* Gère les connexions des clients WebSocket.
* Reçoit les mises à jour de chaque Pokémon.
* Diffuse l'état complet de tous les Pokémon connectés à tous les clients WebSocket et aussi le dashboard.

### Le Serveur gRPC (en Go)

* Gère les flux de communication bidirectionnels avec les clients gRPC.
* Reçoit les mises à jour de chaque Pokémon gRPC.
* Diffuse l'état complet des Pokémon gRPC à tous les clients gRPC connectés.

## Lancer le projet

Le projet utilise Docker et Docker Compose.

1.  Ouvre un terminal à la racine du projet.
2.  Construis les images Docker (si nécessaire) :
    ```bash
    docker-compose build
    ```
3.  Démarre tous les services :
    ```bash
    docker-compose up
    ```
    * Pour lancer plusieurs Pokémon de chaque type, utilise `--scale` :
        ```bash
        # Exemple avec 3 clients WS et 3 clients gRPC
        docker-compose up --build --scale pokemon-client-websocket=3 --scale pokemon-client-grpc=3
        ```

Les services démarrent :
* Serveur gRPC sur le port `50051`.
* Serveur WebSocket sur le port `8080`.

Les logs des clients afficheront la latence RTT moyenne mesurée.

## Latence mesurée
Les logs des clients nous donnent une idée de la latence. Voici un exemple avec un intervalle d'envoi de 1 seconde :

* **Client gRPC** :
    * La latence moyenne s'est montrée très stable, restant autour de 500-510 ms pendant toute la durée du test. Les logs montrent par exemple : `Latence Moyenne RTT = 503.48 ms`, `500.28 ms`, `502.36 ms`.

* **Client WebSocket** :
    * La latence RTT moyenne a démarré à une valeur similaire (environ `523 ms`).
    * Mais apres quelques temps, dans les tests observés, cette latence a augmenté au fil du temps, atteignant (`779 ms`, `2665 ms`, `9144 ms`, `30516 ms`...).

**Interprétation**

* **gRPC** : Dans cette simulation, gRPC a démontré une performance très stable en termes de latence.
* **WebSocket** : La dégradation des performances observée pour WebSocket dans ces logs est notable. Malgres que WebSocket soit généralement performant, cette augmentation de latence suggère un possible problème dans l'implémentation (peut etres la faute du developeur ¯\_(ツ)_/¯ ).

## Quand utiliser quoi ?

* **gRPC:** 
    * Services backend (microservices) doivent communiquer très rapidement entre eux.
    * Performance et une faible latence stable.
    * Contrat d'API fort.
    * Si la communication directe avec un navigateur n'est pas la priorité.

* **WebSocket:** 
    * Développes une application web qui doit afficher des données en temps réel (jeux, dashboards, notifications).
    * Communication bidirectionnelle simple avec le navigateur.
    * Facilité d'intégration côté front-end.
    * (Attention à surveiller les performances  sur la durée, avec l'exemple vu dans ce projet).

## Conclusion

* **gRPC** s'est montré très stable et performant pour le streaming de données dans ce projet de test, idéal pour les communications backend.
* **WebSocket**, bien que parfait pour l'interactivité web, a montré des signes de dégradation de performance dans cette implémentation spécifique, soulignant l'importance des tests et de l'optimisation selon le cas d'usage.

Le choix de la technologie dépendra toujours du contexte spécifique du projet, des besoins en performance et de la facilité d'intégration souhaitée.