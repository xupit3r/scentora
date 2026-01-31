# Specifications Directory

This directory contains detailed technical specifications for the Scentora project.

## Files

### [data-models.md](data-models.md)
Complete database schema and data model specifications:
- Accord, AccordTag, PredefinedTag models
- SQL schema definitions
- Go structs and TypeScript interfaces
- Validation rules and constraints
- Entity relationships

### [api-spec.md](api-spec.md)
REST API endpoint documentation:
- All endpoints with request/response formats
- Authentication and authorization
- Query parameters and filtering
- Error responses and status codes
- Rate limiting specifications

### [tag-system.md](tag-system.md)
Tag system specification:
- 57 predefined tags across 9 categories
- Tag categories (Character, Mood, Season, Time, etc.)
- Custom tag support
- Tag management and filtering
- Database seeding scripts

### [ui-ux-spec.md](ui-ux-spec.md)
User interface and experience design:
- Component specifications (AccordCard, AccordForm, etc.)
- View layouts (Home, Detail, Statistics)
- Color palette and typography
- Interaction patterns and flows
- Responsive design breakpoints
- Accessibility guidelines

## Usage

These specifications serve as the source of truth for implementation. Developers should refer to these documents when:

- Implementing database migrations
- Creating API endpoints
- Building UI components
- Writing tests
- Reviewing code

## Maintenance

Update these specifications when:
- Adding new features
- Modifying data models
- Changing API contracts
- Updating UI designs

Keep specifications in sync with implementation.
