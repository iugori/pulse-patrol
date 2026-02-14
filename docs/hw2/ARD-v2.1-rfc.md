# Architecture Response Document

> **<span style="font-size: 1.3em">Pulse Patrol</span>**
>
> *Develop a software system for healthcare that collects and manages patient data,
> integrates with medical equipment, provides web access for patients and authorized personnel,
> alerts staff for abnormal values, and supports patient transfers between healthcare providers.*

<!-- TOC -->
* [Architecture Response Document](#architecture-response-document)
  * [1. Context](#1-context)
    * [Scope](#scope)
      * [Personas](#personas)
      * [Use cases](#use-cases)
          * [Use Case 1 (Patient)](#use-case-1-patient)
          * [Use Case 2 (Doctor)](#use-case-2-doctor)
          * [Use Case 3 (Doctor)](#use-case-3-doctor)
          * [Use Case 4 (Support Staff)](#use-case-4-support-staff)
          * [Use Case 5 (Administrator)](#use-case-5-administrator)
          * [Use Case 6 (Administrator)](#use-case-6-administrator)
      * [Data Sources and Particularities](#data-sources-and-particularities)
        * [Patient Identification & Demographic Data](#patient-identification--demographic-data)
        * [Electronic Health Records (EHR) & Admission Forms](#electronic-health-records-ehr--admission-forms)
        * [Medical Equipment Telemetry (Real-time Monitoring)](#medical-equipment-telemetry-real-time-monitoring)
        * [Laboratory Information Systems (Test Results)](#laboratory-information-systems-test-results)
        * [Transfer & Continuity of Care Records](#transfer--continuity-of-care-records)
        * [Audit & Access Logs](#audit--access-logs)
    * [Out of Scope](#out-of-scope)
  * [2. Proposed Approach](#2-proposed-approach)
    * [Strategy and Architectural Goals](#strategy-and-architectural-goals)
    * [System Context (C4 Level 1)](#system-context-c4-level-1)
    * [Specialized Microservices](#specialized-microservices)
      * [Care & Clinical Record Service (CCRS)](#care--clinical-record-service-ccrs)
      * [Vital Stream & Alerting Service (VSAS)](#vital-stream--alerting-service-vsas)
      * [Compliance & Identity Service (CIS)](#compliance--identity-service-cis)
      * [Integration & Interop Gateway (IIG)](#integration--interop-gateway-iig)
    * [Data Flow Diagram](#data-flow-diagram)
  * [3. Individual Components Roles and Responsibilities](#3-individual-components-roles-and-responsibilities)
    * [Deployable Units (C4 Level 2)](#deployable-units-c4-level-2)
      * [Use Case Realization](#use-case-realization)
        * [Use Case 1 (Patient)](#use-case-1-patient-1)
        * [Use Case 2 (Doctor)](#use-case-2-doctor-1)
        * [Use Case 3 (Doctor)](#use-case-3-doctor-1)
        * [Use Case 4 (Support Staff)](#use-case-4-support-staff-1)
        * [Use Case 5 (Administrator)](#use-case-5-administrator-1)
        * [Use Case 6 (Administrator)](#use-case-6-administrator-1)
  * [4. Deployment](#4-deployment)
  * [5. Dependencies](#5-dependencies)
  * [6. Data Flows/APIs](#6-data-flowsapis)
    * [Bounded Contexts](#bounded-contexts)
      * [Care Coordination & Admissions](#care-coordination--admissions)
      * [Clinical Records](#clinical-records)
      * [Vital Signs & Monitoring](#vital-signs--monitoring)
      * [Notification & Alerting](#notification--alerting)
      * [Security & Audit (Generic)](#security--audit-generic)
  * [7. Security Concerns](#7-security-concerns)
  * [8. COGS](#8-cogs)
<!-- TOC -->

## 1. Context

### Scope

[//]: # (<<Main system/feature requirements &#40;functional & non-functional&#41;>>)

[//]: # (S: <business-requirements>)

Develop a software system that:

1. Collects data about patients: medical records, test results, admission forms, etc.
2. Collects data from medical equipment used for investigations.
3. The data is accessible to patients through a web application.
4. Doctors and authorized personnel have access to the data of patients admitted to the hospitals where they work.
5. The system allows for alerting medical staff when monitoring systems detect abnormal values.
6. The system can be sold to various healthcare companies and facilitates the transfer of patients from one to another.

[//]: # (S: </business-requirements>)

#### Personas

[//]: # (S: <personas>)

1. **Patient**: Individual receiving medical care or treatment.
    - R3: Patients require access to their health data via a web application.

2. **Doctor**: Medical professional providing care to patients.
    - R4: Doctors require access to the data of patients admitted to the hospitals where they work.
    - R5: They need to receive alerts to respond quickly to patient needs.

3. **Support Staff**: Support personnel assisting in patient care (e.g. nurses).
    - R5: Medical staff need alerts for abnormal values in patient monitoring.

4. **Administrator**: Manager overseeing the healthcare operation.
    - R1: Administrators need to manage patient records effectively.
    - R6: They facilitate patient transfers and ensure proper data handling.

[//]: # (S: </personas>)

#### Use cases

[//]: # (S: <use-cases>)

###### Use Case 1 (Patient)

As a **Patient**,
I want **to access my medical records, test results, and admission forms through a web application**,
so that **I can stay informed about my health status and treatment progress**.

###### Use Case 2 (Doctor)

As a **Doctor**,
I want **to access the data of my patients admitted to the hospital**,
so that **I can provide informed medical care based on their history and current status**.

###### Use Case 3 (Doctor)

As a **Doctor**,
I want **to receive alerts for abnormal values detected by monitoring systems**,
so that **I can respond quickly to critical patient needs and improve outcomes**.

###### Use Case 4 (Support Staff)

As a **Support Staff Member**,
I want **to receive alerts for abnormal values in patient monitoring**,
so that **I can act swiftly to provide necessary medical assistance and ensure patient safety**.

###### Use Case 5 (Administrator)

As an **Administrator**,
I want **to manage patient records effectively**,
so that **I can maintain accurate and up-to-date information for efficient healthcare management**.

###### Use Case 6 (Administrator)

As an **Administrator**,
I want **to facilitate the transfer of patients between healthcare companies**,
so that **I can ensure continuity of care and proper handling of patient data**.

[//]: # (S: </use-cases>)

#### Data Sources and Particularities

[//]: # (S: <data-sources>)

The system manages a complex flow of Protected Health Information (PHI) and Personal Identifiable Information (PII)
that requires a multi-layered approach to legal compliance and technical security.
All data must be encrypted both at rest and in transit, with strict role-based access controls ensuring
that Doctors and Support Staff only access records relevant to their specific hospital.
Under GDPR and HIPAA frameworks, the system must maintain immutable audit logs for every data interaction
to ensure accountability and patient privacy.
Furthermore, to facilitate the seamless transfer of patients between healthcare companies,
the system adheres to data portability standards, providing structured,
machine-readable exports while maintaining the integrity of real-time equipment telemetry used for critical alerting.

##### Patient Identification & Demographic Data

- *Source:* Admission forms and initial registration via the Administrator or Patient (R1).
- *Legal Owner:* The Patient (as the subject) and the Healthcare Provider (as the data controller).

##### Electronic Health Records (EHR) & Admission Forms

- *Source:* Manual entry by Administrators/Support Staff and historical records.
- *Legal Owner:* Usually the Healthcare Provider (Hospital/Clinic) where the data was generated, though patients hold
  rights to access and portability.

##### Medical Equipment Telemetry (Real-time Monitoring)

- *Source:* Direct data streams from medical equipment (e.g., heart rate monitors, ventilators).
- *Legal Owner:* The Healthcare Provider (as the entity conducting the investigation).

##### Laboratory Information Systems (Test Results)

- *Source:* External or internal diagnostic labs.
- *Legal Owner:* The Lab or the Ordering Healthcare Provider.

##### Transfer & Continuity of Care Records

- *Source:* Administrative metadata generated during patient transfers between healthcare companies.
- *Legal Owner:* Joint Ownership or delegated responsibility between the originating and receiving healthcare companies.

##### Audit & Access Logs

- *Source:* System-generated logs of Doctor, Admin, and Patient activity.
- *Legal Owner:* The Healthcare Software Operator/Company.

[//]: # (S: </data-sources>)

### Out of Scope

[//]: # (<<What functional & non-functional requirements we won’t cover in this ARD.>>)

[//]: # (S: <out-of-scope>)

The following items are explicitly excluded from the current architectural design and implementation phase:

- **Medical Equipment Manufacturing/Hardware**: The system integrates with existing equipment via the Integration
  Gateway but does not include the design, maintenance, or manufacturing of the medical hardware itself.

- **Automated Medical Diagnosis**: While the system alerts for "abnormal values" based on predefined thresholds, it will
  not provide automated clinical diagnoses or suggest pharmaceutical treatments (AI-driven medical advice).

- **Direct Billing & Insurance Claims**: Integration with billing systems or processing insurance claims (revenue cycle
  management) is excluded; the focus remains on clinical data and patient transfers.

- **Offline Data Collection**: The system requires an active network connection for real-time alerting. Offline
  buffering and asynchronous syncing from medical devices are not supported in this version.

- **Legal Responsibility for Triage**: The software acts as a communication and monitoring aid. It does not replace the
  professional judgment of medical staff or serve as the primary legal record for emergency dispatch.

- **Identity Provisioning**: The system will integrate with existing Identity Providers (IdP) but will not manage the
  primary creation or physical verification of government-issued identities for patients.

[//]: # (S: </out-of-scope>)

## 2. Proposed Approach

[//]: # (<<How we plan to address the requirements described in the section above. 
The section can include references to ADRs.
It should include a high level C4 context diagram.>>)

The Pulse Patrol architecture is designed to balance the high-availability required for real-time medical alerting with
the strict data isolation and security necessary for multi-tenant healthcare SaaS.

### Strategy and Architectural Goals

- **Microservices Orchestration**: To ensure scalability and independent deployment of the Patient Management and
  Telemetry
  services.

- **Real-time Stream Processing**: Utilizing a pub/sub model (e.g., Kafka or RabbitMQ) to handle high-frequency data
  from
  medical equipment with sub-second latency for alerts.

- **Multi-Tenancy**: A logical isolation strategy allowing the software to be sold to multiple healthcare providers
  while
  ensuring data remains siloed per organization.

- **Interoperability**: Adherence to healthcare standards (like HL7 FHIR) to facilitate seamless patient transfers and
  legacy
  system syncing.

### System Context (C4 Level 1)

The diagram below illustrates how Pulse Patrol sits at the center of the healthcare ecosystem, bridging the gap between
hardware, legacy data, and end-users.

```mermaid
graph LR
    classDef depExt fill: #f8a3a3, stroke: #333, stroke-width: 1.5px;
    classDef depInt fill: #a8e6a1, stroke: #333, stroke-width: 1.5px;
    classDef theSys fill: #92c6ff, stroke: #333, stroke-width: 1.5px;
    classDef groups fill: #f8f8f8, stroke: #333, stroke-width: 1.5px;
    subgraph Legend [Legend]
        direction TB
        L1["External Dependency"]:::depExt
        L2["Internal Dependency"]:::depInt
        L3["Core System"]:::theSys
    end
    subgraph Diagram ["Context Diagram"]
        direction LR
    %% Actors
        subgraph ExternalUsers ["External Users"]
            P(("«person»<br/>👤 Patient&nbsp;")):::depExt
        end

        subgraph InternalUsers ["Internal Users"]
            D(("«person»<br/>👤 Doctor&nbsp;")):::depInt
            S(("«person»<br/>👤 Support&nbsp;<br/>Staff")):::depInt
            A(("«person»<br/>👤 Admin&nbsp;")):::depInt
        end

    %% Core
        PP["«software system»<br/>🫀 Pulse Patrol&nbsp;"]:::theSys
    %% Externals
        subgraph InternalInfrastructure ["Internal Systems"]
            ME["«software system»<br/>📠 Medical Equip.&nbsp;"]:::depInt
            IS["«software system»<br/>💾 Legacy Systems&nbsp;"]:::depInt
        end

        subgraph ExternalInfrastructure ["External Systems"]
            EP["«software system»<br/>🌐 External Peers&nbsp;"]:::depExt
        end

    %% Relationships with concise text
        P -- Accesses records --> PP
        D -- Monitors & Alerts --> PP
        S -- Receives Alerts --> PP
        A -- Manages Transfers --> PP
        PP -- Ingests data from --> ME
        PP <-- Syncs data --> IS
        PP <-- Exchanges data --> EP

    end
%% Styles
    class ExternalUsers groups
    class InternalUsers groups
    class InternalInfrastructure groups
    class ExternalInfrastructure groups
```

### Specialized Microservices

#### Care & Clinical Record Service (CCRS)

Bounded Contexts: Care Coordination & Admissions, Clinical Records.

Supported Functionalities:

- Admission Lifecycle: Managing patient check-ins, registration, and demographic updates.
- Health Record Management: CRUD operations for medical history, allergies, and diagnoses.
- Lab Integration: Processing and storing test results from Laboratory Information Systems (LIS).
- Transfer Orchestration: Packaging patient data into encrypted snapshots for inter-company portability.

Events:

| Name                 | Description                                                            |
|----------------------|------------------------------------------------------------------------|
| PatientAdmitted      | Triggered when a patient is officially checked into a facility.        |
| MedicalRecordUpdated | Emitted when diagnoses, procedures, or allergies are modified.         |
| TestResultPublished  | Produced when laboratory data is imported and ready for viewing.       |
| TransferRequestSent  | Triggered when an Admin initiates a data export to a peer company.     |
| TransferAcknowledged | Produced when a peer company confirms successful receipt of a patient. |
| RecordAccessed       | A privacy event emitted whenever a clinical document is viewed.        |

#### Vital Stream & Alerting Service (VSAS)

Bounded Contexts: Vital Signs & Monitoring, Notification & Alerting.

Supported Functionalities:

- Telemetry Ingestion: High-throughput processing of real-time streams (Heart rate, SpO2, etc.).
- Threshold Evaluation: Real-time analysis of incoming data against predefined "abnormal" triggers.
- Alert Dispatch: Immediate routing of critical notifications to the Clinical Dashboard and mobile devices.
- Short-term Buffer: Maintaining a high-resolution rolling window of physiological data for immediate review.

Events:

| Name                      | Description                                                                    |
|---------------------------|--------------------------------------------------------------------------------|
| TelemetryBatchIngested    | Produced when a window of sensor data (e.g., 5 seconds of SpO2) is stored.     |
| AbnormalThresholdDetected | Emitted the moment the analysis engine identifies a value outside safe limits. |
| AlertTriggered            | Produced when a notification is dispatched to doctors/support staff.           | 
| AlertAcknowledged         | Triggered when a staff member responds to a notification via the dashboard.    |

#### Compliance & Identity Service (CIS)

Bounded Contexts: Security & Audit (Generic).

Supported Functionalities:

- Immutable Auditing: Recording every data access and modification event in a tamper-proof log.
- RBAC & Permissions: Managing "who sees what" based on hospital affiliation and staff role.
- Consent Management: Tracking patient permissions for data sharing and transfers.
- Identity Federation: Bridging with external IdPs for secure staff and patient login.

Events:

| Name                   | Description                                                                          |
|------------------------|--------------------------------------------------------------------------------------|
| UserAuthenticated      | Emitted when a Patient, Doctor, or Admin successfully logs in.                       |
| AccessDenied           | Produced when an unauthorized attempt to view PHI/PII occurs.                        |
| ConsentGranted/Revoked | Triggered when a patient updates their data-sharing preferences.                     |
| AuditLogSealed         | A technical event emitted when a batch of logs is hashed/finalized for immutability. |

#### Integration & Interop Gateway (IIG)

Bounded Contexts: Infrastructure (Cross-cutting).

Supported Functionalities:

- Protocol Translation: Converting legacy messages into internal system events.
- Medical Equipment Adapter: Normalizing various proprietary IoT data formats into a unified stream.
- Peer-to-Peer Handshake: Managing the technical secure tunnel for patient transfers to external healthcare companies.

Events:

| Name                      | Description                                                                            |
|---------------------------|----------------------------------------------------------------------------------------|
| ExternalMessageReceived   | Emitted when a message arrives from a legacy system.                                   |
| EquipmentStreamNormalised | Produced when proprietary IoT data is converted into the standard Pulse Patrol format. |
| TransferRequestReceived   | Triggered when an external peer company attempts to move a patient into the system.    |

### Data Flow Diagram

```mermaid
graph LR
    subgraph Legend [Legend]
        direction LR
        dependency>"Requester"]
        boundary(Sytem<br/>Boundary)
        controller((Controller))
        persistence[(Persistence)]
        queue[[Queue]]
        dependency ~~~ boundary ~~~ controller ~~~ persistence ~~~ queue
    end

    subgraph DataFlow["Data Flow Diagram with Async Event"]
        direction LR
    %% Node Definitions
        P>"«person»<br/>👤 Patient"]
        WP("Web<br/>Portal")
        PMS(("Patient<br/>Management<br/>Services"))
        DS[("Data<br/>Storage")]
        Queue[["Audit"]]
    %% Synchronous Flow
        P -- 1 . Requests<br/>info --> WP
        WP -- 2 . Forwards --> PMS
        PMS -- 3 . Queries --> DS
        DS -- 4 . Returns --> PMS
        PMS -.->|5 . Emit ^RecordAccessed^ event| Queue
        PMS -- 6 . Sends --> WP
        WP -- 7 . Displays --> P

    end

    class CIS async
```

## 3. Individual Components Roles and Responsibilities

[//]: # (<<For each component describe its role and responsibility.
Add container/component and other UML diagrams if needed &#40;sequence&#41;>>)

### Deployable Units (C4 Level 2)

The system will be decomposed into the following functional units:

[//]: # (S: <functional-units>)

- **Web Portal**: Interface for Patients to view records and for Administrators to manage data.
- **Clinical Dashboard**: Specialized interface for Doctors and Support Staff to monitor live telemetry and patient
  data.
- **Patient Management Services**: Core logic for medical records, admission forms, and inter-company transfers.
- **Telemetry & Alerting Services**: Processes real-time data from medical equipment and triggers notifications for
  abnormal values.
- **Data Storage**: Centralized repository for structured medical records and time-series telemetry data.
- **Integration Gateway**: Handles communication with Legacy Systems, Medical Equipment, and Peer Healthcare Companies.

[//]: # (S: </functional-units>)

```mermaid
graph TB
    classDef depExt fill: #f8a3a3, stroke: #333, stroke-width: 1.5px;
    classDef depInt fill: #a8e6a1, stroke: #333, stroke-width: 1.5px;
    classDef theSys fill: white, stroke: #333, stroke-width: 1.5px;
    classDef container fill: #92c6ff, stroke: #333, stroke-width: 1.5px;
    subgraph Legend [Legend]
        direction TB
        L1["External Dependency"]:::depExt
        L2["Internal Dependency"]:::depInt
        L3["Core System"]:::theSys
        L4["Container"]:::container
    end
    Diagram ~~~ Legend
    subgraph Diagram ["Container Diagram"]
        direction TB
        Patient(("«person»<br/>👤 Patient&nbsp;")):::depExt
        Admin(("«person»<br/>👤 Admin&nbsp;")):::depInt
        Doctor(("«person»<br/>👤 Doctor&nbsp;")):::depInt
        Support(("«person»<br/>👤 Support&nbsp;<br/>Staff")):::depInt

        subgraph Pulse_Patrol_System ["«software system» 🫀 Pulse Patrol System Boundary&nbsp;"]
            Portal["«container»<br/>Web Portal&nbsp;"]:::container
            Dashboard["«container»<br/>Clinical Dashboard&nbsp;"]:::container
            PMS["«container»<br/>Patient Management&nbsp;<br/>Services"]:::container
            TAS["«container»<br/>Telemetry & Alerting&nbsp;<br/>Services"]:::container
            Storage[("«container»<br/>Data Storage&nbsp;")]:::container
            Gateway["«container»<br/>Integration Gateway&nbsp;"]:::container
        end

        Peer["«software system»<br/>🌐 Peer Healthcare Companies&nbsp;"]:::depExt
        Equipment["«software system»<br/>📠 Medical Equipment / IoT&nbsp;"]:::depInt
        Legacy["«software system»<br/>💾 Legacy Hospital Systems&nbsp;"]:::depInt
    end
%% User Interactions
    Patient -->|Views records| Portal
    Admin -->|Manages data & transfers| Portal
    Doctor -->|Monitors patients| Dashboard
    Support -->|Receives alerts| Dashboard
%% Internal Routing
    Portal -->|Requests records/transfers| PMS
    Dashboard -->|Fetches live data| TAS
    Dashboard -->|Reads history| PMS
%% Core Logic to Storage
    PMS -->|Read/Write records| Storage
    TAS -->|Store/Retrieve telemetry| Storage
%% Integration Logic
    PMS -->|Syncs data| Gateway
    TAS -->|Receives streams| Gateway
%% External System Connections
    Gateway <-->|Inter - company transfer| Peer
    Equipment -->|Real - time telemetry| Gateway
    Gateway <-->|ETL/Data sync| Legacy
%% Grouping Styling
    class Pulse_Patrol_System theSys
    style Pulse_Patrol_System fill: #33aaff, color: #fff, stroke: #333, stroke-width: 2px

```

#### Use Case Realization

##### Use Case 1 (Patient)

As a **Patient**,
I want **to access my medical records, test results, and admission forms through a web application**,
so that **I can stay informed about my health status and treatment progress**.

```mermaid
sequenceDiagram
    actor P as «person»<br/>Patient<br/>
    participant WP as «container»<br/>Web Portal
    participant PMS as «container»<br/>Patient Management Services
    participant DS as «container»<br/>Data Storage
    P ->> WP: Request access to medical records
    WP ->> PMS: Forward request for patient data
    PMS ->> DS: Retrieve medical records
    DS -->> PMS: Send medical data
    PMS -->> WP: Provide patient data
    WP -->> P: Display medical records

```

##### Use Case 2 (Doctor)

As a **Doctor**,
I want **to access the data of my patients admitted to the hospital**,
so that **I can provide informed medical care based on their history and current status**.

```mermaid
sequenceDiagram
    actor D as «person»<br/>Doctor<br/>
    participant CD as «container»<br/>Clinical Dashboard
    participant PMS as «container»<br/>Patient Management Services
    participant TAS as «container»<br/>Telemetry & Alerting Services
    participant DS as «container»<br/>Data Storage
    D ->> CD: Select Patient Profile
%% Fetching Medical History
    CD ->> PMS: Get Patient Records (History, Labs, Admissions)
    PMS ->> DS: Query Patient History
    DS -->> PMS: Return Structured Records
    PMS -->> CD: Patient Medical History
%% Fetching Real-time Status
    CD ->> TAS: Get Live Telemetry Data
    TAS ->> DS: Query Recent Sensor Data
    DS -->> TAS: Return Time-series Data
    TAS -->> CD: Current Physiological Status
    CD -->> D: Display Integrated Patient View
```

##### Use Case 3 (Doctor)

As a **Doctor**,
I want **to receive alerts for abnormal values detected by monitoring systems**,
so that **I can respond quickly to critical patient needs and improve outcomes**.

```mermaid
sequenceDiagram
    participant ME as «software system»<br/>📠 Medical Equipment
    participant IG as «container»<br/>Integration Gateway
    participant TAS as «container»<br/>Telemetry & Alerting Services
    participant DS as «container»<br/>Data Storage
    participant CD as «container»<br/>Clinical Dashboard
    actor SS as «person»<br/>Doctor<br/>
    Note over ME, SS: Real-time Monitoring & Alerting Flow
    ME ->> IG: Stream real-time telemetry data
    IG ->> TAS: Forward data stream

    rect rgb(240, 240, 240)
        Note right of TAS: Process & Analyze Values
        TAS ->> TAS: Detect abnormal threshold breach
    end

    par Concurrent Actions
        TAS ->> DS: Store telemetry & alert record
        TAS ->> CD: Push critical alert notification
    end

    CD -->> SS: Visual/Audible Alert: "Abnormal Values Detected"
    SS ->> CD: Acknowledge alert & view patient details
    CD ->> DS: Fetch latest patient vitals
    DS -->> CD: Return historical & current data
    CD -->> SS: Display comprehensive patient status
```

##### Use Case 4 (Support Staff)

As a **Support Staff Member**,
I want **to receive alerts for abnormal values in patient monitoring**,
so that **I can act swiftly to provide necessary medical assistance and ensure patient safety**.

```mermaid
sequenceDiagram
    participant ME as «software system»<br/>📠 Medical Equipment
    participant IG as «container»<br/>Integration Gateway
    participant TAS as «container»<br/>Telemetry & Alerting Services
    participant DS as «container»<br/>Data Storage
    participant CD as «container»<br/>Clinical Dashboard
    actor D as «person»<br/>Support Staff<br/>
    Note over ME, D: Real-time Patient Monitoring Flow
    ME ->> IG: Stream Telemetry Data (e.g., Heart Rate, SpO2)
    IG ->> TAS: Forward Raw Data Stream
    TAS ->> TAS: Process & Analyze Values

    alt Abnormal Values Detected
        TAS ->> DS: Log Alert Event & Telemetry
        TAS -->> CD: Push Real-time Critical Alert
        CD -->> D: Visual/Audible Alert Notification
        D ->> CD: Acknowledge Alert & View Live Feed
        CD ->> TAS: Request Detailed Telemetry
        TAS ->> CD: Stream High-Resolution Data
    else Normal Values
        TAS ->> DS: Store Routine Telemetry
    end
    Note right of D: Doctor intervenes based on alert
```

##### Use Case 5 (Administrator)

As an **Administrator**,
I want **to manage patient records effectively**,
so that **I can maintain accurate and up-to-date information for efficient healthcare management**.

```mermaid
sequenceDiagram
    actor Admin as «person»<br/>Admin<br/>
    participant Portal as «container»<br/>Web Portal
    participant PMS as «container»<br/>Patient Management Services
    participant Storage as «container»<br/>Data Storage
    participant Gateway as «container»<br/>Integration Gateway
    participant Legacy as «software system»<br/>📠 Legacy Hospital Systems
    Admin ->> Portal: Access Patient Record
    Portal ->> PMS: Request Patient Data
    PMS ->> Storage: Fetch Record
    Storage -->> PMS: Record Data
    PMS -->> Portal: Display Record
    Admin ->> Portal: Submit Record Updates
    Portal ->> PMS: UpdatePatient(Data)

    rect rgb(240, 240, 240)
        Note over PMS, Legacy: Atomic Transaction
        PMS ->> Storage: Persist Updated Record
        Storage -->> PMS: Confirmation
        PMS ->> Gateway: Sync Update to Legacy
        Gateway ->> Legacy: Update EMR/Legacy DB
        Legacy -->> Gateway: Sync Success
        Gateway -->> PMS: Sync Acknowledged
    end

    PMS -->> Portal: Update Successful
    Portal -->> Admin: Display Success Message
```

##### Use Case 6 (Administrator)

As an **Administrator**,
I want **to facilitate the transfer of patients between healthcare companies**,
so that **I can ensure continuity of care and proper handling of patient data**.

```mermaid
sequenceDiagram
    actor Admin as «person»<br/>Admin<br/>
    participant Portal as «container»<br/>Web Portal
    participant PMS as «container»<br/>Patient Management Services
    participant Storage as «container»<br/>Data Storage
    participant Gateway as «container»<br/>Integration Gateway
    participant Peer as «software system»<br/>🌐 Peer Healthcare Company
    Note over Admin, Peer: Use Case: Transfer Patient to Another Company
    Admin ->> Portal: Select patient and target healthcare provider
    Portal ->> PMS: InitiateTransfer(patientId, targetCompanyId)
    PMS ->> Storage: FetchPatientFullRecord(patientId)
    Storage -->> PMS: Medical records, test results, history
    PMS ->> PMS: Package and Encrypt Patient Data
    PMS ->> Gateway: SendTransferRequest(targetCompanyId, encryptedData)
    Gateway ->> Peer: PostPatientTransfer(encryptedData)

    alt Transfer Successful
        Peer -->> Gateway: 201 Created (Transfer Acknowledgement)
        Gateway -->> PMS: Success Status
        PMS ->> Storage: UpdatePatientStatus(Transferred / Archived)
        PMS -->> Portal: Transfer confirmed
        Portal -->> Admin: Display "Transfer Completed Successfully"
    else Transfer Failed
        Peer -->> Gateway: Error (e.g., Validation Failed)
        Gateway -->> PMS: Error Status
        PMS -->> Portal: Transfer Failed Notification
        Portal -->> Admin: Display "Transfer Error - Please retry"
    end
```

## 4. Deployment

[//]: # (<<Include a deployment diagram and documentation about it - regions, communication, networking, etc.>>)

## 5. Dependencies

[//]: # (<<Both internal and external dependencies.
For example, Plane&Simple has an external dependency of a payment system>>)

Entities interacting with the Pulse Patrol system:

[//]: # (S: <external-entities>)

- Human Actors
    - *External*
        - **Patients** - views personal medical history, test results, and treatment progress via the web portal
    - *Internal*
        - **Doctor** - accesses patient data within their hospital and receives critical physiological alerts
        - **Support Staff Member** - nurses/assistants who receive real-time alerts for abnormal patient monitoring
          values
        - **Administrator** - manages records, oversees data integrity, and initiates inter-company patient transfers
- Technical Systems
    - *External*
        - **External Healthcare Companies Peer** - systems belonging to other providers that receive or send patient
          data during a transfer
    - *Internal*
        - **Medical Equipment** - IoT devices and monitoring hardware (e.g., bedside monitors, ventilators) that stream
          real-time telemetry
        - **Legacy Hospital Systems** - existing legacy databases or EMRs where admission forms and historical medical
          records may reside

[//]: # (S: </external-entities>)

## 6. Data Flows/APIs

[//]: # (<<Data flow diagrams. Definitions of component’s APIs>>)

### Bounded Contexts

#### Care Coordination & Admissions

This context manages the lifecycle of a patient’s presence within a healthcare facility and the handover of
responsibility between organizations. It acts as the system’s "entry and exit" gatekeeper, ensuring that every patient
is correctly identified and that their care journey remains continuous when moving between peer providers.

*Events*

- *PatientAdmitted*: An individual is officially registered in a hospital facility.
- *TransferRequestSent*: Data export was triggered for an external provider.
- *TransferRequestReceived*: Data export was triggered from an external provider.
- *TransferAcknowledged*: The receiving company confirmed receipt of patient data (triggered).

*Entities*

- Patient - the main patient record (identification, demographics, etc.)
- Peer - information about a peer healthcare facility
- MedicalSnapshot - a bundle of information sent / received to / from peer (possibly binary data)

*Aggregates*

- Admission - information for each patient admission
- Transfer - information for each patient transfer

```mermaid
graph TD
    classDef aggregate rx: 50, ry: 50, fill: LightSkyBlue, stroke: #333, stroke-width: 1.5px;
    classDef boundedContext rx: 100, ry: 50, color: white, fill: MidnightBlue, stroke: #333, stroke-width: 1.5px;
    classDef entityValue fill: LightGreen, stroke: #333, stroke-width: 1.5px;
    subgraph Legend [Legend]
        direction TB
        L1["Domain"]:::boundedContext
        L2["Aggregate"]:::aggregate
        L3["Entity / Value Object"]:::entityValue
    end
    CCA ~~~ Legend

    subgraph CCA ["Care Coordination & Admissions"]
        subgraph Admission ["«aggregate»"]
            direction TB
            AdmissionRoot["«root»<br/>Admission"]:::entityValue
            AdmissionPatient["«valueObject»<br/>PatientID"]:::entityValue
            AdmissionDate["«valueObject»<br/>AdmissionDate"]:::entityValue
            AdmissionRoot ---> AdmissionPatient
            AdmissionRoot --> AdmissionDate
        end

        subgraph Transfer ["«aggregate»"]
            direction TB
            TransferRoot["«root»<br/>Transfer"]:::entityValue
            TransferPatient["«valueObject»<br/>PatientID"]:::entityValue
            TransferPeer["«valueObject»<br/>PeerID"]:::entityValue
            TransferMedicalSnapshot["«valueObject»<br/>MedicalSnapshotID"]:::entityValue
            TransferDate["«valueObject»<br/>TransferDate"]:::entityValue
            TransferRoot ---> TransferPatient
            TransferRoot ---> TransferPeer
            TransferRoot ---> TransferMedicalSnapshot
            TransferRoot --> TransferDate
        end

        Patient["«entity»<br/>Patient"]:::entityValue
        Peer["«entity»<br/>Peer"]:::entityValue
        MedicalSnapshot["«entity»<br/>MedicalSnapshot"]:::entityValue
        Admission ~~~ Patient
        Transfer ~~~ Peer
        Transfer ~~~ MedicalSnapshot
    end
%% Styling to make the subgraph look like an oval/capsule
    class CCA boundedContext
    class Admission aggregate
    class Transfer aggregate
```

#### Clinical Records

This context serves as the authoritative source of truth for a patient's medical history. It manages the lifecycle and
integrity of static clinical data, ensuring that both patients and providers have a consistent view of health progress.

*Events*

- *MedicalRecordUpdated*: Changes to clinical history were successfully persisted.
- *TestResultPublished*: Laboratory results were made available for viewing.
- *RecordAccessed*: An authorized person viewed a specific clinical document.

*Entities*

- PatientProfile - a local representation of the patient’s clinical identity
- StaffMemberProfile - a local representation for each concerned staff member
- MedicalHistoryEntry - individual entries such as diagnoses, procedures, or allergies
- TestResult - the actual data (e.g., "Glucose: 100 mg/dL").
- ReportMetadata - details about the performing lab, timestamps, and the ordering physician

*Aggregates*

- ClinicalRecord - ensures that all medical documentation is tied to a specific patient and hospital context
- LabResult - handles the specific complexities of diagnostic data coming from the Laboratory Information Systems

```mermaid
graph TD
    classDef aggregate rx: 50, ry: 50, fill: LightSkyBlue, stroke: #333, stroke-width: 1.5px;
    classDef boundedContext rx: 100, ry: 50, color: white, fill: MidnightBlue, stroke: #333, stroke-width: 1.5px;
    classDef entityValue fill: LightGreen, stroke: #333, stroke-width: 1.5px;
    subgraph Legend [Legend]
        direction TB
        L1["Domain"]:::boundedContext
        L2["Aggregate"]:::aggregate
        L3["Entity / Value Object"]:::entityValue
    end
    CR ~~~ Legend

    subgraph CR ["Clinical Records"]
        direction TB

        subgraph ClinicalRecord ["«aggregate»"]
            direction TB
            ClinicalRecordRoot["«root»<br/>ClinicalRecord"]:::entityValue
            ClinicalRecordPatient["«valueObject»<br/>PatientID"]:::entityValue
            ClinicalRecordHistoryEntry["«valueObject»<br/>MedicalHistoryEntryID"]:::entityValue
            ClinicalRecordRoot --> ClinicalRecordPatient
            ClinicalRecordRoot --> ClinicalRecordHistoryEntry
            LabResultRoot["«aggregate»<br/>LabResult"]:::aggregate
            LabResultTestResult["«valueObject»<br/>TestResultID"]:::entityValue
            LabResultReportMetadata["«valueObject»<br/>ReportMetadataID"]:::entityValue
            LabResultStaffMemberProfile["«valueObject»<br/>StaffMemberProfileID"]:::entityValue
            LabResultRoot --> LabResultTestResult
            LabResultRoot --> LabResultReportMetadata
            LabResultRoot --> LabResultStaffMemberProfile
            ClinicalRecordRoot ---> LabResultRoot
        end

        PatientProfile["«entity»<br/>PatientProfile"]:::entityValue
        MedicalHistoryEntry["«entity»<br/>MedicalHistoryEntry"]:::entityValue
        StaffMemberProfile["«entity»<br/>StaffMemberProfile"]:::entityValue
        TestResult["«entity»<br/>TestResult"]:::entityValue
        ReportMetadata["«entity»<br/>ReportMetadata"]:::entityValue
    %% Layout anchors
        ClinicalRecord ~~~ PatientProfile
        ClinicalRecord ~~~ TestResult
        ClinicalRecord ~~~ StaffMemberProfile
        ClinicalRecord ~~~ MedicalHistoryEntry
        ClinicalRecord ~~~ ReportMetadata
    end

    class CR boundedContext
    class ClinicalRecord aggregate
```

#### Vital Signs & Monitoring

Description This context is responsible for the continuous ingestion and evaluation of real-time physiological data
(telemetry) from medical devices. It acts as the system’s "nervous system," observing incoming streams to identify
critical changes in a patient’s state. Its primary role is to distinguish between normal physiological patterns and
urgent clinical deviations or technical failures.

*Events*

- *AbnormalValueDetected*: A vital sign breached a predefined safety threshold.
- *EquipmentDisconnected*: The data stream from the device was lost.

*Entities*

- MedicalDevice - represents the physical hardware (e.g., Bedside Monitor ID)
- VitalSignReading - a single data point (value, unit, timestamp)
- ThresholdConfig the defined "safe" ranges for a specific patient

*Aggregates*

- MonitoringSession - binds a Patient to a MedicalDevice for a specific duration

```mermaid
graph TD
    classDef aggregate rx: 50, ry: 50, fill: LightSkyBlue, stroke: #333, stroke-width: 1.5px;
    classDef boundedContext rx: 100, ry: 50, color: white, fill: MidnightBlue, stroke: #333, stroke-width: 1.5px;
    classDef entityValue fill: LightGreen, stroke: #333, stroke-width: 1.5px;
    subgraph Legend [Legend]
        direction TB
        L1["Domain"]:::boundedContext
        L2["Aggregate"]:::aggregate
        L3["Entity / Value Object"]:::entityValue
    end
    VSM ~~~ Legend

    subgraph VSM ["Vital Signs & Monitoring"]
        direction TB

        subgraph MonitoringSession ["«aggregate»"]
            direction TB
            MonitoringSessionRoot["«root»<br/>MonitoringSession"]:::entityValue
            MonitoringSessionMedicalDevice["«valueObject»<br/>MedicalDeviceID"]:::entityValue
            MonitoringSessionVitalSignReading["«valueObject»<br/>VitalSignReadingID"]:::entityValue
            MonitoringSessionThresholdConfig["«valueObject»<br/>ThresholdConfigID"]:::entityValue
            MonitoringSessionRoot --> MonitoringSessionMedicalDevice
            MonitoringSessionRoot --> MonitoringSessionVitalSignReading
            MonitoringSessionRoot --> MonitoringSessionThresholdConfig
        end

        MedicalDevice["«entity»<br/>MedicalDevice"]:::entityValue
        VitalSignReading["«entity»<br/>VitalSignReading"]:::entityValue
        ThresholdConfig["«entity»<br/>ThresholdConfig"]:::entityValue
    %% Layout anchors
        MonitoringSession ~~~ MedicalDevice
        MonitoringSession ~~~ VitalSignReading
        MonitoringSession ~~~ ThresholdConfig
    end

    class VSM boundedContext
    class MonitoringSession aggregate
```

#### Notification & Alerting

This context is responsible for the lifecycle of a notification, from the moment a telemetry threshold is breached to
the final acknowledgment by a human operator. It decouples the detection of a medical issue from the delivery of the
message.

*Events*

- *AlertTriggered*: A notification was created based on abnormal vitals.
- *StaffNotified*: The alert was successfully delivered to a device.
- *AlertAcknowledged*: A medical professional responded to the notification.

*Entities*

- Recipient - a projection of the Staff member (from the Identity/Staff context) containing their active device tokens
  and availability status
- NotificationChannel - represents the medium used to reach staff (e.g., Push Notification, SMS, Dashboard Popup)

*Aggregates*

- Alert - the central record of a specific abnormal event. It tracks the severity, the source (Patient/Device), and the
  current state (Triggered, Notified, Acknowledged)

```mermaid
graph TD
    classDef aggregate rx: 50, ry: 50, fill: LightSkyBlue, stroke: #333, stroke-width: 1.5px;
    classDef boundedContext rx: 100, ry: 50, color: white, fill: MidnightBlue, stroke: #333, stroke-width: 1.5px;
    classDef entityValue fill: LightGreen, stroke: #333, stroke-width: 1.5px;
    subgraph Legend [Legend]
        direction TB
        L1["Domain"]:::boundedContext
        L2["Aggregate"]:::aggregate
        L3["Entity / Value Object"]:::entityValue
    end
    NA ~~~ Legend

    subgraph NA ["Notification & Alerting"]
        direction TB

        subgraph Alert ["«aggregate»"]
            direction TB
            AlertRoot["«root»<br/>Alert"]:::entityValue
            AlertRecipient["«valueObject»<br/>RecipientID"]:::entityValue
            AlertNotificationChannel["«valueObject»<br/>NotificationChannelID"]:::entityValue
            AlertMessage["«valueObject»<br/>Message"]:::entityValue
            AlertRoot ---> AlertRecipient
            AlertRoot ---> AlertNotificationChannel
            AlertRoot --> AlertMessage
        end

        Recipient["«entity»<br/>Recipient"]:::entityValue
        NotificationChannel["«entity»<br/>NotificationChannel"]:::entityValue
    %% Layout anchors
        Alert ~~~ Recipient
        Alert ~~~ NotificationChannel
    end

    class NA boundedContext
    class Alert aggregate
```

#### Security & Audit (Generic)

<u>
This section is added as a reminder.
More details will be added after security module will be researched in detail.
</u>

*Events*

- *AccessGranted*: Permission was successfully verified for a data request.
- *UnauthorizedAccessDetected*: A security breach attempt was recognized.
- ...

## 7. Security Concerns

[//]: # (<<Authorisation, Authentication, Data encryption, Threat modelling diagram>>)

## 8. COGS

[//]: # (<<Cost estimation model for hardware, services, data storage and transfer for the whole solution>>)
