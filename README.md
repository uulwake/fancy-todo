# Fancy Todo

Implement fancy todo API with 5 different languages such as Typescript, Go, Java, Rust, and Elixir.

## Product Requirement

### User

1. User can register with name, email and password
2. User can login with email and password

### Task

1. User can create task with title, description, and tags
2. User can get detail of the task
3. User can get list of tasks
4. User can filter tasks based on status and tag
5. User can sort tasks based on specific field with specific order
6. User can adjust the number of tasks in one page
7. User can update title, description, status, and order of a task
8. User can delete a task
9. User can search task by title

### Tag

1. User can create tag
2. User can add existing tag to a task
3. User can delete a tag
4. User can search tag by name

## API Contract V1

### User

1. Register User

   - Endpoint: `POST /v1/users/register`
   - Request:
     ```json
     {
       "name": "foobar",
       "email": "foobar@mail.com",
       "password": "secret"
     }
     ```
   - Response:
     ```json
     {
       "data": {
         "user": {
           "id": 1
         },
         "jwt_token": "token"
       }
     }
     ```

2. Login User
   - Endpoint: `POST /v1/users/login`
   - Request:
     ```json
     {
       "email": "foobar@mail.com",
       "password": "secret"
     }
     ```
   - Response:
     ```json
     {
       "data": {
         "user": {
           "id": 1
         },
         "jwt_token": "token"
       }
     }
     ```

### Task

1. Create Task

   - Endpoint: `POST /v1/tasks`
   - Request:
     ```json
     {
       "title": "Title 1",
       "description": "Desc 1",
       "tag_ids": [1, 2, 3]
     }
     ```
   - Response:
     ```json
     {
       "data": {
         "task": {
           "id": 1
         }
       }
     }
     ```

2. Get Detail Task
   - Endpoint: `GET /v1/tasks/:task_id`
   - Response:
     ```json
     {
       "data": {
         "task": {
           "id": 1,
           "title": "Title 1",
           "description": "Desc 1",
           "status": "on_going",
           "order": 1,
           "created_at": "2023-12-12 12:12:12",
           "updated_at": "2023-12-12 12:12:12"
         }
       }
     }
     ```
3. Get List of Tasks
   - Endpoint: `GET /v1/tasks`
   - Query:
     - status: string. Ex: `?status=completed`
     - page_size: number. Ex: `?page_size=10`
     - page_number: number. Ex: `?page_number=2`
     - tag_id: number, Ex: `tag_id=1`
     - sort_name: string, Ex: `sort_key=title`
     - sort_order: string, Ex: `sort_order=desc`
   - Response:
     ```json
     {
       "data": {
         "tasks": [
           {
             "id": 1,
             "title": "Title 1",
             "status": "on_going",
             "order": 1,
             "created_at": "2023-12-12 12:12:12",
             "updated_at": "2023-12-12 12:12:12"
           }
         ]
       },
       "page": {
         "size": 5,
         "number": 1,
         "total": 10
       }
     }
     ```
4. Update Task - Reorder, Title, Content, Tag, Status
   - Endpoint: `PATCH /v1/tasks/:task_id`
   - Request:
     ```json
     {
       "title": "Title 2",
       "description": "Desc 2",
       "status": "completed",
       "order": 2
     }
     ```
   - Response:
     ```json
     {
       "data": {
         "task": {
           "id": 1
         }
       }
     }
     ```
5. Delete Task

   - Endpoint: `DELETE /v1/tasks/:task_id`
   - Response:
     ```json
     {
       "data": {
         "task": {
           "id": 1
         }
       }
     }
     ```

6. Search Task by Title
   - Endpoint: `GET /v1/tasks/search`
   - Query:
     - title: string. Ex: `?title=Titl`
   - Response:
     ```json
     {
       "data": {
         "tasks": [
           {
             "id": 1,
             "user_id": 2,
             "title": "Title 1",
             "status": "on_going"
           }
         ]
       }
     }
     ```

### Tag

1. Add New Tag in Post

   - Endpoint: `POST /v1/tags`
   - Request:
     ```json
     {
       "name": "tag_name",
       "task_id": 1
     }
     ```
   - Response
     ```json
     {
       "data": {
         "tag": {
           "id": 1
         }
       }
     }
     ```

2. Add Existing Tag in Task

   - Endpoint: `PATCH /v1/tags/:tag_id/tasks/:task_id`
   - Request:
     ```json
     {}
     ```
   - Response:
     ```json
     {
       "data": {
         "tag": {
           "id": 1,
           "task": {
             "id": 2
           }
         }
       }
     }
     ```

3. Delete Tag

   - Endpoint: `DELETE /v1/tags/tag_id`
   - Response:
     ```json
     {
       "data": {
         "tag": {
           "id": 1
         }
       }
     }
     ```

4. Search Tag by Name

   - Endpoint: `GET /v1/tags/search`
   - Query
     - name: string. Ex: `?name=tag_`
   - Response:

   ```json
   {
     "data": {
       "tags": [
         {
           "id": 1,
           "name": "tag_1"
         },
         {
           "id": 2,
           "name": "tag_2"
         }
       ]
     }
   }
   ```
