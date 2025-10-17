# Taxpayer Signup Implementation Guide

## Overview
The seamless taxpayer signup flow has been successfully implemented. This allows regular users (taxpayers) to register with both authentication credentials and taxpayer profile information in a single request.

## Key Features
- **One-to-One Relationship**: Links users and taxpayers via `user_id` foreign key
- **Role-Based Access**: Uses `role = 'user'` to distinguish taxpayers from admin/employees
- **Centralized Auth**: Keeps authentication in `users` table, profile data in `taxpayers` table
- **Compliance Ready**: Sensitive data like `national_id` stays in taxpayers table

## API Endpoint

### POST /auth/register
Register a new taxpayer with both user credentials and taxpayer profile.

**Request Body:**
```json
{
  "email": "john.doe@example.com",
  "password": "securepassword123",
  "first_name": "John",
  "last_name": "Doe",
  "phone_number": "+254712345678",
  "county_id": 1,
  "role": "user",
  "taxpayer_type": "individual",
  "national_id": "12345678",
  "business_name": ""
}
```

**For Business Taxpayers:**
```json
{
  "email": "business@example.com",
  "password": "securepassword123",
  "first_name": "Jane",
  "last_name": "Smith",
  "phone_number": "+254712345679",
  "county_id": 1,
  "role": "user",
  "taxpayer_type": "business",
  "national_id": "87654321",
  "business_name": "Smith Enterprises Ltd"
}
```

**Response:**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "email": "john.doe@example.com",
  "first_name": "John",
  "last_name": "Doe",
  "role": "user"
}
```

## Validation Rules

### Required Fields (All Users)
- `email`: Valid email address
- `password`: Minimum 8 characters
- `first_name`: User's first name
- `last_name`: User's last name
- `phone_number`: Contact phone number
- `county_id`: Valid county ID

### Taxpayer-Specific Fields (role = "user")
- `taxpayer_type`: Must be "individual" or "business"
- `national_id`: Unique national identification number
- `business_name`: Required only for business taxpayers

## Database Schema

### Users Table
```sql
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    county_id INTEGER REFERENCES counties(id),
    email VARCHAR UNIQUE NOT NULL,
    password_hash VARCHAR NOT NULL,
    first_name VARCHAR NOT NULL,
    last_name VARCHAR NOT NULL,
    phone_number VARCHAR,
    role VARCHAR NOT NULL DEFAULT 'user',
    employee_id VARCHAR,
    department VARCHAR,
    is_active BOOLEAN DEFAULT true,
    last_login TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### Taxpayers Table
```sql
CREATE TABLE taxpayers (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    county_id INTEGER NOT NULL REFERENCES counties(id),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    taxpayer_type VARCHAR NOT NULL CHECK (taxpayer_type IN ('individual', 'business')),
    national_id VARCHAR UNIQUE NOT NULL,
    email VARCHAR NOT NULL,
    phone_number VARCHAR,
    first_name VARCHAR,
    last_name VARCHAR,
    business_name VARCHAR,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

## Implementation Details

### Auth Service Enhancement
The existing `auth.RegisterRequest` has been extended with taxpayer-specific fields:
- `taxpayer_type`: Individual or business classification
- `national_id`: Unique taxpayer identifier
- `business_name`: Business name for business taxpayers

### Automatic Profile Creation
When a user registers with `role = "user"`, the system automatically:
1. Creates a user record in the `users` table
2. Creates a linked taxpayer profile in the `taxpayers` table
3. Validates taxpayer-specific fields
4. Ensures national_id uniqueness

### Error Handling
- Validates all required fields before processing
- Checks for existing users with the same email
- Checks for existing taxpayers with the same national_id
- Returns appropriate error messages for validation failures

## Usage Examples

### Test with cURL

**Individual Taxpayer Registration:**
```bash
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john.doe@example.com",
    "password": "securepassword123",
    "first_name": "John",
    "last_name": "Doe",
    "phone_number": "+254712345678",
    "county_id": 1,
    "role": "user",
    "taxpayer_type": "individual",
    "national_id": "12345678"
  }'
```

**Business Taxpayer Registration:**
```bash
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "business@example.com",
    "password": "securepassword123",
    "first_name": "Jane",
    "last_name": "Smith",
    "phone_number": "+254712345679",
    "county_id": 1,
    "role": "user",
    "taxpayer_type": "business",
    "national_id": "87654321",
    "business_name": "Smith Enterprises Ltd"
  }'
```

### Login After Registration
After successful registration, users can login using the existing login endpoint:

```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john.doe@example.com",
    "password": "securepassword123"
  }'
```

## Next Steps

### Post-Signup Flow
After successful registration and login, taxpayers can:
1. Access their tax dashboard
2. View tax obligations
3. File tax returns
4. Make payments
5. View payment history

### Additional Queries
The implementation includes a `GetFullProfileByUserID` query that joins user and taxpayer data for complete profile information:

```sql
SELECT u.id, u.email, u.first_name, u.last_name, u.role,
       t.county_id, t.taxpayer_type, t.national_id, t.business_name
FROM users u
JOIN taxpayers t ON u.id = t.user_id
WHERE u.id = $1 AND u.is_active = true;
```

This enables efficient retrieval of complete taxpayer profiles for dashboard and profile management features.

## Security Considerations
- Passwords are hashed using bcrypt
- National IDs are stored separately from auth credentials
- User accounts can be deactivated without losing taxpayer data
- Foreign key constraints ensure data integrity
- Sensitive taxpayer data is isolated from general user data