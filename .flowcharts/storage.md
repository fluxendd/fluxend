### What is storage?
A backend system needs storage. Your app may have profile pictures, covers, and other files. You need to store them somewhere. You can use a cloud storage service like AWS S3, Google Cloud Storage, or Azure Blob Storage. You can also use a local file system or a database to store files. Fluxton handles this for you.

You have several supported storage drivers to choose from. These are:

- Local
- S3
- Dropbox
- Backblaze

### How storage works?
You can configure driver of your choice. Fluxton takes care of rest. The result endpoints for listing, creating, updating, deleting containers and files are the same for all drivers. You can use any driver you want. 

```mermaid
graph LR
    A[File upload request received] --> B[Resolve user storage config]
    B --> C[Load selected storage driver]
    C --> D{Driver type}

    D --> E[S3 Service]
    D --> F[Dropbox Service]
    D --> G[Backblaze Service]
    D --> H[Google Drive Service]
    D --> I[OneDrive Service]
    D --> J[Wasabi Service]
    D --> K[Local Storage Handler]

    E --> L[Upload to S3 API]
    F --> M[Upload to Dropbox API]
    G --> N[Upload to Backblaze API]
    H --> O[Upload to Google Drive API]
    I --> P[Upload to OneDrive API]
    J --> Q[Upload to Wasabi API]
    K --> R[Write file to disk]

    L --> Z[Upload complete]
    M --> Z
    N --> Z
    O --> Z
    P --> Z
    Q --> Z
    R --> Z

```
