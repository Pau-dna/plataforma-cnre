# Endpoint de Progreso Reciente del Usuario

## Descripción

Este endpoint permite obtener los 10 registros de progreso más recientes de un usuario específico de la tabla `user_progress`, incluyendo las relaciones precargadas con las tablas `modules` y `contents`.

## Endpoint

```
GET /api/v1/users/{userId}/recent-progress
```

### Parámetros

- `userId` (path parameter, requerido): ID del usuario del cual se quiere obtener el progreso reciente

### Respuesta Exitosa (200)

Retorna un array de objetos `UserProgress` con las relaciones `Module` y `Content` precargadas:

```json
[
  {
    "id": 1,
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T10:30:00Z",
    "user_id": 123,
    "course_id": 1,
    "module_id": 2,
    "content_id": 5,
    "completed_at": "2024-01-15T10:30:00Z",
    "score": 95,
    "attempts": 1,
    "module": {
      "id": 2,
      "title": "Módulo de Ejemplo",
      "description": "Descripción del módulo",
      "order": 1,
      "course_id": 1
    },
    "content": {
      "id": 5,
      "title": "Lección de Ejemplo",
      "description": "Descripción de la lección",
      "type": "content",
      "order": 1,
      "module_id": 2
    }
  }
]
```

### Respuestas de Error

- **400 Bad Request**: ID de usuario inválido
- **500 Internal Server Error**: Error al obtener el progreso reciente del usuario

## Implementación

### Backend (Go)

#### Repository Layer
```go
func (r *userprogressRepository) GetRecentByUser(userID uint, limit int) ([]*models.UserProgress, error)
```
- Consulta la base de datos filtrando por `user_id`
- Ordena por `created_at DESC` para obtener los más recientes
- Limita los resultados al número especificado (10)
- Precarga las relaciones `Module` y `Content` usando GORM Preload

#### Service Layer
```go
func (s *userProgressService) GetRecentUserProgress(userID uint) ([]*models.UserProgress, error)
```
- Llama al método del repositorio con un límite de 10 registros
- Maneja errores con mensajes en español

#### Handler Layer
```go
func (h *UserProgressHandler) GetRecentUserProgress(c *gin.Context)
```
- Valida el parámetro `userId`
- Llama al servicio
- Retorna la respuesta JSON con los registros de progreso

### Frontend (TypeScript)

#### Controller
```typescript
async getRecentUserProgress(userId: number): Promise<UserProgress[]>
```
- Método del controlador para consumir el endpoint desde el frontend
- Retorna un array de objetos `UserProgress`

## Uso

### Desde el Frontend
```typescript
import { UserProgressController } from '$lib/controllers';

const controller = new UserProgressController();
const recentProgress = await controller.getRecentUserProgress(userId);

// Acceder a los datos precargados
recentProgress.forEach(progress => {
  console.log(`Módulo: ${progress.module?.title}`);
  console.log(`Contenido: ${progress.content?.title}`);
  console.log(`Completado: ${progress.completed_at}`);
});
```

### Desde una API externa
```bash
curl -X GET "http://localhost:8080/api/v1/users/123/recent-progress" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

## Notas

- El endpoint retorna exactamente 10 registros o menos si el usuario tiene menos de 10 progresos
- Los registros se ordenan por fecha de creación descendente (más reciente primero)
- Las relaciones `Module` y `Content` siempre se precargan para evitar consultas N+1
- Los mensajes de error están en español según las convenciones del proyecto
