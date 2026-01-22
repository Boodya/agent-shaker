# Dashboard Counters Implementation

## Overview
Created a comprehensive dashboard with real-time statistics counters from the backend API, split into reusable Vue components.

## Backend Changes

### 1. New Dashboard Handler (`internal/handlers/dashboard.go`)
Created a new handler that provides comprehensive statistics:

**Endpoint:** `GET /api/dashboard`

**Response Structure:**
```json
{
  "projects": {
    "total": 10,
    "active": 8,
    "archived": 2
  },
  "agents": {
    "total": 15,
    "active": 12,
    "idle": 2,
    "offline": 1
  },
  "tasks": {
    "total": 45,
    "pending": 10,
    "in_progress": 15,
    "done": 18,
    "blocked": 2
  },
  "contexts": {
    "total": 30
  }
}
```

**Features:**
- Uses PostgreSQL `COUNT(*) FILTER (WHERE ...)` for efficient aggregation
- Error handling with fallback to zero values
- Separate statistics for each entity type

### 2. Updated Server Configuration (`cmd/server/main.go`)
- Added `dashboardHandler` initialization
- Registered `/api/dashboard` route
- Handler connected to database for real-time data

## Frontend Changes

### 1. New StatCard Component (`web/src/components/StatCard.vue`)
**Purpose:** Reusable statistics card with breakdown details

**Props:**
- `title` (String, required) - Card title
- `value` (Number/String, required) - Main statistic value
- `icon` (String) - Emoji icon
- `iconBgColor` (String) - Background color for icon
- `breakdown` (Array) - Detailed breakdown with labels, values, and colors

**Features:**
- Hover animations and transitions
- Flexible breakdown section
- Color-coded values
- Responsive design

**Example Usage:**
```vue
<StatCard
  title="Projects"
  :value="stats.projects.total"
  icon="📁"
  iconBgColor="#3b82f6"
  :breakdown="[
    { label: 'Active', value: stats.projects.active, color: '#10b981' },
    { label: 'Archived', value: stats.projects.archived, color: '#6b7280' }
  ]"
/>
```

### 2. Updated Dashboard View (`web/src/views/Dashboard.vue`)
**New Features:**
- Fetches dashboard stats from `/api/dashboard` endpoint
- Uses StatCard component for all statistics
- Loading and error states
- Four main stat cards:
  1. **Projects** - Total, Active, Archived
  2. **Agents** - Total, Active, Idle, Offline
  3. **Tasks** - Total, Pending, In Progress, Blocked
  4. **Completed** - Done tasks + Contexts count

**Components Structure:**
- StatCard components for metrics
- Recent Projects section (unchanged)
- Active Agents section (unchanged)
- Recent Tasks section (unchanged)

## Statistics Breakdown

### Projects Counter
- **Total:** All projects in the system
- **Active:** Projects with status = 'active'
- **Archived:** Projects with status = 'archived'

### Agents Counter
- **Total:** All registered agents
- **Active:** Agents currently working (status = 'active')
- **Idle:** Agents available but not working (status = 'idle')
- **Offline:** Agents not connected (status = 'offline')

### Tasks Counter
- **Total:** All tasks created
- **Pending:** Tasks waiting to start (status = 'pending')
- **In Progress:** Tasks currently being worked on (status = 'in_progress')
- **Blocked:** Tasks that are blocked (status = 'blocked')

### Completed Counter
- **Done:** Successfully completed tasks (status = 'done')
- **Contexts:** Total documentation/context entries

## Benefits

1. **Real-time Data:** Statistics fetched directly from database
2. **Reusable Component:** StatCard can be used elsewhere in the application
3. **Visual Hierarchy:** Color-coded breakdown for quick insights
4. **Scalable:** Easy to add more statistics
5. **Performance:** Single API call for all dashboard data
6. **Maintainable:** Separated concerns between backend and frontend

## API Usage Example

```bash
# Fetch dashboard statistics
curl http://localhost:8080/api/dashboard

# Response
{
  "projects": {"total": 5, "active": 4, "archived": 1},
  "agents": {"total": 8, "active": 6, "idle": 1, "offline": 1},
  "tasks": {"total": 20, "pending": 5, "in_progress": 8, "done": 6, "blocked": 1},
  "contexts": {"total": 12}
}
```

## Future Enhancements

1. Add time-based filtering (today, this week, this month)
2. Add trend indicators (↑ ↓) comparing to previous period
3. Add charts/graphs for visual representation
4. Add real-time updates via WebSocket
5. Add export functionality for reports
6. Add user-specific statistics
7. Add performance metrics (avg completion time, etc.)

## Testing

1. **Backend:** Build verified with `go build ./cmd/server`
2. **Frontend:** No compilation errors in Vue components
3. **API:** Endpoint registered and handler connected to database

## Files Modified

**Backend:**
- `internal/handlers/dashboard.go` (NEW)
- `cmd/server/main.go` (Modified - added dashboard handler and route)

**Frontend:**
- `web/src/components/StatCard.vue` (NEW)
- `web/src/views/Dashboard.vue` (Modified - integrated StatCard and API call)
