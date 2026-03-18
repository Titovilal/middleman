# Claude The Middleman (CTM)

## Concepto

Un orquestador de agentes Claude donde el chat principal actúa como **gestor puro**: nunca toca archivos ni hace trabajo técnico directamente. Su único rol es decidir qué agente trabaja, cuándo, y con qué historial.

---

## Problema que resuelve

Cuando trabajas en proyectos complejos con Claude Code, terminas con una lista de instancias abiertas:
- Cada una tiene un contexto diferente (archivos leídos, decisiones tomadas, errores vistos)
- Tienes que recordar mentalmente "este Claude sabe del módulo X", "este tiene el contexto del refactor Y"
- Haces rewinds manuales cuando un cambio reciente contamina el contexto de una instancia
- Decides tú mismo a qué instancia mandar cada pregunta

CTM automatiza y externaliza esa carga cognitiva.

---

## Rol del Middleman (chat principal)

El Middleman **nunca**:
- Lee archivos del proyecto
- Ejecuta comandos
- Hace cambios de código

El Middleman **solo**:
- Mantiene un registro de agentes activos y su contexto conocido
- Decide a qué agente delegar cada tarea
- Ordena rewinds cuando detecta que un cambio reciente perjudica el contexto
- Recupera historial de conversaciones anteriores para restaurar un agente a un estado útil
- Fusiona o descarta agentes según su utilidad actual

---

## Estructura de agentes

Cada agente tiene:
- **ID / nombre**: identificador semántico (ej: `agent-auth`, `agent-db-refactor`)
- **Contexto conocido**: qué archivos ha leído, qué decisiones se tomaron en su sesión
- **Checkpoint**: punto del historial al que se puede hacer rewind
- **Estado**: activo / en espera / descartado

---

## Flujo de trabajo

```
Usuario → Middleman → decide agente adecuado → delega tarea
                    ↓
              si el agente tiene contexto contaminado → rewind a checkpoint
                    ↓
              si ningún agente tiene el contexto → crear agente nuevo con briefing
                    ↓
              agente trabaja internamente (el Middleman no ve el stream)
                    ↓
              agente devuelve solo la respuesta final + metadata de contexto
                    ↓
              Middleman actualiza registro → responde al usuario
```

### Principio de caja negra

El Middleman **nunca ve** el flujo interno de un agente: herramientas usadas, archivos leídos, razonamiento intermedio. Solo recibe:
- La respuesta final
- Metadata opcional: qué archivos tocó, si hubo errores, si el contexto se amplió

Esto es deliberado: mantiene el contexto del Middleman limpio y enfocado en orquestación, no en detalles de implementación.

---

## Operaciones del Middleman

| Operación | Descripción |
|---|---|
| `spawn(name, briefing)` | Crea un agente nuevo con contexto inicial |
| `delegate(agent_id, task)` | Envía una tarea al agente más adecuado |
| `rewind(agent_id, checkpoint)` | Restaura un agente a un punto anterior del historial |
| `inspect(agent_id)` | Consulta el estado y contexto conocido de un agente |
| `merge(agent_a, agent_b)` | Combina el conocimiento de dos agentes en uno |
| `discard(agent_id)` | Marca un agente como descartado |
| `history(agent_id)` | Recupera el historial de decisiones y cambios de un agente |

---

## Por qué tiene sentido

- El contexto es el recurso escaso. Contaminar un agente con una dirección equivocada es costoso.
- El rewind no es un error — es una estrategia deliberada para preservar contexto limpio.
- Tener un gestor que no consume contexto técnico significa que siempre tiene ancho de banda para razonar sobre la orquestación.
- Separa claramente "pensar qué hacer" de "hacer".

---

## Implementación posible

- El Middleman es una instancia de Claude Code con un sistema prompt que le prohíbe usar herramientas de archivo/código
- Cada agente es una instancia separada de Claude Code (subagent / worktree / sesión independiente)
- El registro de agentes se persiste en un archivo JSON o markdown en el proyecto
- Los checkpoints son commits git o snapshots del historial de conversación exportados

---

## Estado actual

Idea inicial. Pendiente de diseño técnico detallado.
