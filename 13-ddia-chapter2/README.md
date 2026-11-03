### 1. Comparison Table: Relational vs. Document vs. Graph

| Feature | Relational Model | Document Model | Graph Model |
| :--- | :--- | :--- | :--- |
| **Primary Structure** | Tables and rows | Self-contained JSON documents | Vertices and directed edges |
| **Schema Enforcement** | Schema-on-write | Schema-on-read | Schema-flexible |
| **Best Relationship Fit** | Routine many-to-many joins | One-to-many nested trees | Complex variable-depth webs |
| **Storage Locality** | Split across tables | Contiguous per document | Optimized for hop traversal |
| **Query Language** | SQL | Declarative JSON or API | Cypher or SPARQL |

*The above table compares the three models across schema behavior, relationship structure, disk locality, and query interfaces.*

### 2. Social Network Feature Sketch

```text
+-------------------------------------------------------------------+
|                   Social Network Go API Backend                   |
+-----------------+-------------------------------+-----------------+
                  |                               |
1. Fetch/Create Post            2. Add/Count Likes       3. Friend Suggestions
     (One-to-Many Tree)              (High Write Volume)      (Multi-Hop Traversal)
                  |                               |                 |
                  v                               v                 v
    +---------------------------+   +---------------------------+   +------------------------+
    |      Document Store       |   |      Relational DB        |   |       Graph DB         |
    |  {                        |   |  Table: post_likes        |   |  (UserA)-[:FRIENDS]->  |
    |    "post_id": 901,        |   |  +---------+-----------+  |   |  (UserB)-[:FRIENDS]->  |
    |    "author_id": 42,       |   |  | post_id | user_id   |  |   |  (UserC)-[:ATTENDED]-> |
    |    "text": "Hello",       |   |  +---------+-----------+  |   |  (University: AASTU)   |
    |    "media": [...]         |   |  | 901     | 104       |  |   |                        |
    |  }                        |   |  +---------+-----------+  |   |                        |
    +---------------------------+   +---------------------------+   +------------------------+

    ### 3. Architectural Justification

A modern social network splits storage across three specialized data models to match each feature's read and write access patterns.

First, **Posts** fit the Document model. A post is a self-contained tree containing text, media metadata, and poll options. Storing each post as a single JSON document provides strong storage locality, loading the entire post in one disk read without multi-table joins, while schema-on-read accommodates new post attachments without lock-heavy table migrations.

Second, **Likes** fit the Relational model. Appending likes to an array inside a post document causes severe disk reallocation and write-lock contention when a post goes viral. A normalized relational table with a composite primary key (`post_id`, `user_id`) enforces one-like-per-user uniqueness at the database level, handles thousands of concurrent tiny inserts safely, and executes fast `COUNT(*)` aggregations.

Third, **Friendships** fit the Graph model. Social connections form a many-to-many web where friend recommendations require variable-depth traversals across friends of friends, shared universities, and cities. While relational databases struggle with expensive recursive joins across unknown hop counts, a property graph traverses multi-hop relationship chains rapidly and evolves cleanly as new entity types connect.