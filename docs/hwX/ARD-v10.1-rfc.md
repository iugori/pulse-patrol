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
      * [Container relationships diagram](#container-relationships-diagram)
      * [Container communication diagram](#container-communication-diagram)
      * [Use Case Realization](#use-case-realization)
        * [Sequence 1 (Patient)](#sequence-1-patient)
        * [Sequence 2 (Doctor)](#sequence-2-doctor)
        * [Sequence 3 (Doctor)](#sequence-3-doctor)
        * [Sequence 4 (Support Staff)](#sequence-4-support-staff)
        * [Sequence 5 (Administrator)](#sequence-5-administrator)
        * [Sequence 6 (Administrator)](#sequence-6-administrator)
  * [4. Deployment](#4-deployment)
    * [Runtime Technologies](#runtime-technologies)
      * [Compute & Frontend (Nodes)](#compute--frontend-nodes)
      * [Data Storage & Caching](#data-storage--caching)
      * [Communication (Edges)](#communication-edges)
    * [Development & Maintenance Technologies](#development--maintenance-technologies)
    * [AWS Overview](#aws-overview)
      * [Consolidated Technology Summary](#consolidated-technology-summary)
      * [SWOT Analysis: AWS AI Technologies](#swot-analysis-aws-ai-technologies)
    * [Deployment Diagram](#deployment-diagram)
      * [Key Architectural Components](#key-architectural-components)
  * [5. Dependencies](#5-dependencies)
  * [6. Data Flows/APIs](#6-data-flowsapis)
    * [Bounded Contexts](#bounded-contexts)
      * [Care Coordination & Admissions](#care-coordination--admissions)
      * [Clinical Records](#clinical-records)
      * [Vital Signs & Monitoring](#vital-signs--monitoring)
      * [Notification & Alerting](#notification--alerting)
      * [Security & Audit (Generic)](#security--audit-generic)
    * [API Specifications](#api-specifications)
      * [Patient Transfer API](#patient-transfer-api)
  * [7. Security Concerns](#7-security-concerns)
    * [Data Flow Diagram for  Use Case 6 (Patient Transfer)](#data-flow-diagram-for--use-case-6-patient-transfer)
      * [Asset Identification Table](#asset-identification-table)
      * [Threat Identification Table](#threat-identification-table)
      * [Security Controls & Mitigation Table](#security-controls--mitigation-table)
    * [Compliance with CIA Principles (Confidentiality, Integrity, Availability)](#compliance-with-cia-principles-confidentiality-integrity-availability)
  * [8. COGS](#8-cogs)
    * [8.1 Regions for Comparison](#81-regions-for-comparison)
    * [8.2 Selected AWS Services for Estimation](#82-selected-aws-services-for-estimation)
    * [8.3 Estimation Assumptions](#83-estimation-assumptions)
    * [8.4 Results](#84-results)
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
- **Compliance & Identity Services**: Compliance & Identity Services: Handles authentication, authorization, audit...

[//]: # (S: </functional-units>)

##### Container relationships diagram

[//]: # (S: <container-diagram-relationships>)

This diagram shows the logical relationships between containers, focusing on data dependencies.

```mermaid
graph TB
    classDef depExt fill: #f8a3a3, stroke: #333, stroke-width: 1.5px;
    classDef depInt fill: #a8e6a1, stroke: #333, stroke-width: 1.5px;
    classDef theSys fill: white, stroke: #333, stroke-width: 1.5px;
    classDef container fill: #92c6ff, stroke: #333, stroke-width: 1.5px;
%%-
    subgraph Legend [Legend]
        direction TB
        L1["External Dependency"]:::depExt
        L2["Internal Dependency"]:::depInt
        L3["Core System"]:::theSys
        L4["Container"]:::container
    end
    Diagram ~~~ Legend
%%-
    subgraph Diagram ["Container Diagram"]
        direction TB
        Patient(("«person»<br/>👤 Patient&nbsp;")):::depExt
        Admin(("«person»<br/>👤 Admin&nbsp;")):::depInt
        Doctor(("«person»<br/>👤 Doctor&nbsp;")):::depInt
        Support(("«person»<br/>👤 Support&nbsp;<br/>Staff")):::depInt

        subgraph Pulse_Patrol_System ["«software system» 🫀 Pulse Patrol System Boundary&nbsp;"]
            Portal["«container»<br/>Web Portal<br/>[S3 + CloudFront]"]:::container
            Dashboard["«container»<br/>Clinical Dashboard<br/>[S3 + CloudFront]"]:::container
            PMS["«container»<br/>Patient Management<br/>[ECS/Fargate]"]:::container
            TAS["«container»<br/>Telemetry & Alerting<br/>[Lambda/Kinesis]"]:::container
            Storage[("«container»<br/>Data Storage<br/>[Aurora + Timestream]")]:::container
            Gateway["«container»<br/>Integration Gateway<br/>[IoT Core / API Gateway]"]:::container
            CIS["«container»<br/>Compliance & Identity<br/>[ECS/Fargate]"]:::container
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

[//]: # (S: </container-diagram-relationships>)

##### Container communication diagram

This diagram shows the technical communication protocols and patterns (sync arrows, async dotted arrows with queues).

[//]: # (S: <container-diagram-communication>)

```mermaid
graph TB
    classDef depExt fill: #f8a3a3, stroke: #333, stroke-width: 1.5px;
    classDef depInt fill: #a8e6a1, stroke: #333, stroke-width: 1.5px;
    classDef theSys fill: white, stroke: #333, stroke-width: 1.5px;
    classDef container fill: #92c6ff, stroke: #333, stroke-width: 1.5px;
%%-
    subgraph Legend [Legend]
        direction TB
        L1["External Dependency"]:::depExt
        L2["Internal Dependency"]:::depInt
        L3["Core System"]:::theSys
        L4["Container"]:::container
        Source["Source ⦅ initiates communication ⦆"] -->|Sync| Target
        Source["Source ⦅ initiates communication ⦆"] -.->|Async| Target
    end
    Diagram ~~~ Legend
%%-
    subgraph Diagram ["Container Diagram (Inter-Service Communication)"]
        direction LR
        Patient(("<br/>«person»<br/><span style='font-size:50px'>👤</span><br/><br/>&nbsp;&nbsp;Patient&nbsp;&nbsp;")):::depExt
        Admin(("<br/>«person»<br/><span style='font-size:50px'>👤</span><br/><br/>&nbsp;&nbsp;&nbsp;Admin&nbsp;&nbsp;&nbsp;")):::depInt
        Doctor(("<br/>«person»<br/><span style='font-size:50px'>👤</span><br/><br/>&nbsp;&nbsp;&nbsp;Doctor&nbsp;&nbsp;&nbsp;")):::depInt
        Support(("«person»<br/><br/><span style='font-size:56px'>👤</span><br/><br/>Support Staff")):::depInt

        subgraph Pulse_Patrol_System ["«software system» 🫀 Pulse Patrol System Boundary"]
            direction TB
            Portal["«container»<br/>Web Portal (uiWP)"]:::container
            Dashboard["«container»<br/>Clinical Dashboard (uiCD)"]:::container
            Gateway["«container»<br/>Integration Gateway (sGW)"]:::container

            subgraph ssPMS["&nbsp;"]
                PMSq[["«queue»<br/>PM"]]:::container
                PMS["«container»<br/>Patient Management (sPM)"]:::container
                Spms[("«container»<br/>Data Storage (pPM)")]:::container
            end

            subgraph ssTAS["&nbsp;"]
                TASq[["«queue»<br/>TA"]]:::container
                TAS["«container»<br/>Telemetry & Alerting (sTA)"]:::container
                Stas[("«container»<br/>Data Storage (pTA)")]:::container
            end

            subgraph ssAAA["&nbsp;"]
                direction LR
                AAAq[["«queue»<br/>Audit"]]:::container
                AAA["«container»<br/>Compliance & Identity (sAAA)"]:::container
                Saaa[("«container»<br/>Data Storage (pAAA)")]:::container
            end
        end

        Peer["«software system»<br/>🌐 Peer Healthcare (ePEER)"]:::depExt
        Equipment["«software system»<br/>📠 Medical Equipment (iEQP)"]:::depInt
        Legacy["«software system»<br/>💾 Legacy Systems (iLEG)"]:::depInt
    end

%% User to UI
    Patient -->|HTTPS| Portal
    Admin -->|HTTPS| Portal
    Doctor <-->|HTTPS + WSS| Dashboard
    Support <-->|HTTPS + WSS| Dashboard
%% UI to Services
    Portal -->|OIDC/OAuth2| AAA
    Dashboard -->|OIDC/OAuth2| AAA
    Portal -->|REST| PMS
    Dashboard <-->|REST/gRPC| TAS
    Dashboard -->|REST/gRPC| PMS
    PMS -->|⦅ Read ⦆| PMSq
    TAS -->|⦅ Read ⦆| TASq
%% Integration & Telemetry
    Equipment -->|MQTT QoS 0/2| Gateway
    Gateway <-->|HL7 FHIR / MLLP| Legacy
    Gateway <-->|mTLS REST| Peer
    Gateway -..->|Async Message| PMSq
    Gateway -.->|Async Message| TASq
%% Backend to Persistence & Audit
    PMS -->|DB Driver| Spms
    TAS -->|DB Driver| Stas
    AAA -->|DB Driver| Saaa
    PMS -.->|Async Audit Queue| AAAq
    TAS -.->|Async Audit Queue| AAAq
    Gateway -.->|Async Audit Queue| AAAq
    AAA -->|⦅ Read ⦆| AAAq
%% Reactive Alerts
    TAS -->|Mobile Push| Doctor
    TAS -->|Mobile Push| Support
%% Styling
    class Pulse_Patrol_System theSys
    style Pulse_Patrol_System fill: #f5f5f5, stroke: #333, stroke-width: 2px
```

[//]: # (S: </container-diagram-communication>)

#### Use Case Realization

[//]: # (S: <use-case-1>)

##### Sequence 1 (Patient)

As a **Patient**,
I want **to access my medical records, test results, and admission forms through a web application**,
so that **I can stay informed about my health status and treatment progress**.

```mermaid
sequenceDiagram
    autonumber
    actor P as «person»<br/>Patient<br/>
    participant WP as «container»<br/>Web Portal<br/>(uiWP)
    participant AAA as «container»<br/>Compliance & Identity<br/>Services (sAAA)
    participant PMS as «container»<br/>Patient Management<br/>Services (sPM)
    participant DS as «container»<br/>Data Storage<br/>(pDS)
    Note over P, AAA: Identity & Auth (OIDC/OAuth2 over HTTPS)
    P ->> WP: Access records
    activate WP
    WP ->> AAA: Authenticate User & Validate Token
    activate AAA
    AAA -->> WP: Token Validated / Scopes Authorized
    deactivate AAA
    Note over WP, PMS: Synchronous Request (REST/HTTPS)
    WP ->> PMS: GET /patient-records (with Identity Token)
    activate PMS

    par Asynchronous Audit Logging
        PMS -) AAA: Emit Audit Event (Actor, Patient, Timestamp)
        Note right of AAA: Non-blocking Fire-and-Forget
    and Synchronous Data Retrieval
        PMS ->> DS: Query medical_records schema
        activate DS
        DS -->> PMS: Return records/results
        deactivate DS
    end

    PMS -->> WP: 200 OK (JSON Data)
    deactivate PMS
    WP -->> P: Display medical history & test results
    deactivate WP
```

[//]: # (S: </use-case-1>)

##### Sequence 2 (Doctor)

As a **Doctor**,
I want **to access the data of my patients admitted to the hospital**,
so that **I can provide informed medical care based on their history and current status**.

```mermaid
sequenceDiagram
    autonumber
    actor D as «person»<br/>Doctor<br/>
    participant CD as «container»<br/>Clinical Dashboard<br/>(uiCD)
    participant AAA as «container»<br/>Compliance & Identity<br/>Services (sAAA)
    participant PMS as «container»<br/>Patient Management<br/>Services (sPM)
    participant TAS as «container»<br/>Telemetry & Alerting<br/>Services (sTA)
    participant DS as «container»<br/>Data Storage<br/>(pDS)
    D ->> CD: Select Patient Profile
    activate CD
    CD ->> AAA: Verify Session/Permissions (OIDC)
    activate AAA
    AAA -->> CD: Identity Token (JWT)
    deactivate AAA

    rect rgb(240, 248, 255)
        Note over CD, PMS: Synchronous REST: Historical Data
        CD ->> PMS: GET /patients/{id}/records
        activate PMS
        par Async Audit
            PMS -) AAA: Log Access Event
        and Data Fetch
            PMS ->> DS: Query Relational Schema
            activate DS
            DS -->> PMS: Historical Records
            deactivate DS
        end
        PMS -->> CD: Patient Medical History
        deactivate PMS
    end

    rect rgb(255, 245, 238)
        Note over CD, TAS: Reactive WSS: Live Telemetry
        CD ->> TAS: Establish Stream / Fetch Config
        activate TAS
        TAS ->> DS: Query Time-Series Schema
        activate DS
        DS -->> TAS: Recent Telemetry Data
        deactivate DS
        TAS -->> CD: Push Live Telemetry (WSS)
        deactivate TAS
    end

    CD -->> D: Display Integrated Patient View
    deactivate CD
```

##### Sequence 3 (Doctor)

As a **Doctor**,
I want **to receive alerts for abnormal values detected by monitoring systems**,
so that **I can respond quickly to critical patient needs and improve outcomes**.

```mermaid
sequenceDiagram
    autonumber
    participant ME as «software system»<br/>📠 Medical Equipment<br/>(iEQP)
    participant IG as «container»<br/>Integration Gateway<br/>(sGW)
    participant TAS as «container»<br/>Telemetry & Alerting<br/>Services (sTA)
    participant DS as «container»<br/>Data Storage<br/>(pDS)
    participant CD as «container»<br/>Clinical Dashboard<br/>(uiCD)
    participant AAA as «container»<br/>Compliance & Identity<br/>Services (sAAA)
    actor D as «person»<br/>Doctor<br/>
    Note over ME, IG: Protocol: MQTT (QoS 2 - Exactly Once)
    ME ->> IG: Stream Telemetry/Alert Packet
    activate IG
    IG -) TAS: Push Normalized Data<br/>(Async Message Bus)
    deactivate IG
    activate TAS

    rect rgb(255, 230, 230)
        Note right of TAS: Analysis: Detect Threshold Breach
        TAS ->> TAS: Trigger "Critical Alert" Logic
    end

    par Async Persistence & Audit
        TAS ->> DS: Store Time-Series Data & Alert
        TAS -) AAA: Emit Audit Log (Abnormal Event)
    and Reactive Notification
        TAS -->> CD: Push Alert (WebSockets / Mobile Push)
        activate CD
    end
    deactivate TAS
    CD -->> D: Visual/Audible Alert: "Vitals Critical!"
    D ->> CD: Acknowledge & View Details
    CD ->> TAS: Fetch Focused Real-Time Stream
    activate TAS
    TAS -->> CD: Stream Live Vitals
    deactivate TAS
    CD -->> D: Display Live Patient Vitals
    deactivate CD
```

##### Sequence 4 (Support Staff)

As a **Support Staff Member**,
I want **to receive alerts for abnormal values in patient monitoring**,
so that **I can act swiftly to provide necessary medical assistance and ensure patient safety**.

```mermaid
sequenceDiagram
    autonumber
    participant ME as «software system»<br/>📠 Medical Equipment<br/>(iEQP)
    participant IG as «container»<br/>Integration Gateway<br/>(sGW)
    participant TAS as «container»<br/>Telemetry & Alerting<br/>Services (sTA)
    participant AAA as «container»<br/>Compliance & Identity<br/>Services (sAAA)
    participant DS as «container»<br/>Data Storage<br/>(pDS)
    participant CD as «container»<br/>Clinical Dashboard<br/>(uiCD)
    actor SSM as «person»<br/>Support Staff<br/>
    Note over ME, IG: Stream: MQTT (QoS 0 for Waveforms)
    ME ->> IG: Push Raw Telemetry
    activate IG
    IG -) TAS: Async Message Broker (Normalized Data)
    deactivate IG
    activate TAS

    rect rgb(255, 240, 245)
        Note right of TAS: Analysis Engine
        TAS ->> TAS: Detect Threshold Breach
    end

    alt Abnormal Values Detected
        par Non-Blocking Operations
            TAS ->> DS: Store Alert & Time-Series Data
            TAS -) AAA: Fire-and-Forget Audit Event
        and Reactive Push
            TAS -->> CD: Push Alert (WebSockets)
            activate CD
            CD -->> SSM: Visual/Audible Alert Notification
        end

        SSM ->> CD: Acknowledge & Open Live Feed
        CD ->> TAS: Request High-Res Stream
        TAS -->> CD: Stream High-Resolution Data
        deactivate TAS
    else Normal Values
        TAS ->> DS: Batch Store Routine Telemetry
    end
    deactivate CD
```

##### Sequence 5 (Administrator)

As an **Administrator**,
I want **to manage patient records effectively**,
so that **I can maintain accurate and up-to-date information for efficient healthcare management**.

```mermaid
sequenceDiagram
    autonumber
    actor A as «person»<br/>Administrator<br/>
    participant WP as «container»<br/>Web Portal<br/>(uiWP)
    participant AAA as «container»<br/>Compliance & Identity<br/>Services (sAAA)
    participant PMS as «container»<br/>Patient Management<br/>Services (sPM)
    participant DS as «container»<br/>Data Storage<br/>(pDS)
    participant GW as «container»<br/>Integration Gateway<br/>(sGW)
    participant LEG as «software system»<br/>💾 Legacy Hospital Systems<br/>(iLEG)
    A ->> WP: Access Patient Record
    activate WP
    WP ->> AAA: Authenticate & Check Admin Role (OIDC)
    activate AAA
    AAA -->> WP: Identity Token (JWT)
    deactivate AAA
    WP ->> PMS: GET /patients/{id}
    activate PMS
    PMS ->> DS: Fetch Record (Relational Schema)
    activate DS
    DS -->> PMS: Record Data
    deactivate DS
    PMS -->> WP: Display Record
    deactivate PMS
    A ->> WP: Submit Updated Records/Forms
    WP ->> PMS: POST /patients/{id}/update (with Token)
    activate PMS
    Note over PMS, DS: Update Internal Source of Truth
    PMS ->> DS: Persist Updated Record
    activate DS
    DS -->> PMS: Confirmation
    deactivate DS

    par Async Audit & Legacy Sync
        PMS -) AAA: Emit Audit Log (Admin Update Event)
        Note over PMS, GW: Modern Standards (JSON/REST)
        PMS ->> GW: Notify Record Change
        activate GW
        Note over GW, LEG: Protocol Translation (HL7/Flat-file)
        GW ->> LEG: Update Legacy EMR/DB
        activate LEG
        LEG -->> GW: Sync Success
        deactivate LEG
        GW -->> PMS: Sync Acknowledged
        deactivate GW
    end

    PMS -->> WP: Update Successful
    deactivate PMS
    WP -->> A: Display Success Message
    deactivate WP
```

##### Sequence 6 (Administrator)

As an **Administrator**,
I want **to facilitate the transfer of patients between healthcare companies**,
so that **I can ensure continuity of care and proper handling of patient data**.

```mermaid
sequenceDiagram
    autonumber
    actor A as «person»<br/>Administrator<br/>
    participant WP as «container»<br/>Web Portal<br/>(uiWP)
    participant AAA as «container»<br/>Compliance & Identity<br/>Services (sAAA)
    participant PMS as «container»<br/>Patient Management<br/>Services (sPM)
    participant DS as «container»<br/>Data Storage<br/>(pDS)
    participant GW as «container»<br/>Integration Gateway<br/>(sGW)
    participant PEER as «software system»<br/>🌐 External Peer<br/>(ePEER)
%% PHASE 1: VALIDATION
    A ->> WP: Initiate patient transfer & select peer
    activate WP
    WP ->> AAA: Verify Admin credentials (OIDC)
    activate AAA
    AAA -->> WP: Identity Token (JWT)
    deactivate AAA
    WP ->> PMS: POST /api/v1/transfers (patientId, targetProviderId, ...)
    activate PMS
    PMS ->> DS: 1.1 Verify Consent & Target Provider
    activate DS
    DS -->> PMS: Consent Validated
    deactivate DS
%% PHASE 2: DATA AGGREGATION
    PMS ->> DS: 2.1 Fetch Records (Aurora) & Vitals (Timestream)
    activate DS
    DS -->> PMS: Full Patient Data Set
    deactivate DS
%% PHASE 3: SECURITY (PACKAGING)
    rect rgb(230, 255, 230)
        Note over PMS: 3.1 Encrypt PII (AES-256-GCM)<br/>3.2 Generate SHA-256 Hash<br/>3.3 Sign with Private Key
    end
%% ASYNC HANDOFF
    PMS -->> WP: 202 Accepted (transferId & trackingUrl)
    deactivate PMS
    WP -->> A: Display "Transfer Processing (ID: 9f3d...)"
    deactivate WP
%% PHASE 4: TRANSFER (BACKGROUND)
    Note over PMS, PEER: Background Workflow Commences
    activate PMS
    PMS ->> GW: 4.1 Forward Encrypted Payload
    activate GW
    Note over GW, PEER: Protocol: mTLS REST API (Handshake)
    GW ->> PEER: 4.2 POST /incoming-transfer
    activate PEER

    alt Transfer Successful
        PEER -->> GW: 201 Created (Acknowledgement)
        GW -->> PMS: Transfer Confirmed
    %% PHASE 5: FINALIZATION
        par Finalization Tasks
            PMS ->> DS: 5.1 Mark record as "Transferred"
            PMS -) AAA: 5.2 Emit Audit Event
        end
    else Transfer Failed
        PEER -->> GW: 5xx Error / Timeout
        deactivate PEER
        GW -->> PMS: 4.3 Report Failure
        deactivate GW
        PMS -) AAA: 5.3 Log Failed Transfer Attempt
    end
    deactivate PMS
```

Basis for [Patient Transfer API](#patient-transfer-api).

## 4. Deployment

[//]: # (<<Include a deployment diagram and documentation about it - regions, communication, networking, etc.>>)

### Runtime Technologies

This section maps the architectural nodes and edges to live AWS infrastructure.

#### Compute & Frontend (Nodes)

* **uiWP & uiCD (Web/Clinical Dashboards):** Hosted on **Amazon S3** (Static Website) and distributed via **Amazon
  CloudFront**.
* **sPM, sTA, sAAA (Microservices):**
  * *Primary:* **AWS Fargate (on Amazon ECS)**. Serverless container orchestration
    provides scaling without managing VMs.
  * *Alternative:* **AWS Lambda**. Lower cost for low-traffic services, but potential "cold starts" are risky for
    sTA (
    Telemetry).
  * **Trade-off:** Fargate is better for the 24/7 uptime required by hospitals.


* **sGW (Integration Gateway):** **AWS IoT Core** for MQTT (iEQP) and **Amazon API Gateway** for REST.

#### Data Storage & Caching

* **Structured Records (sPM):** **Amazon Aurora (PostgreSQL)**. HIPAA-compliant and supports relational integrity.
* **High-Frequency Telemetry (sTA):** **Amazon Timestream**. Optimized for time-series data from equipment.
* **Clinical Reports/Images (S3):** **Amazon S3** with **S3 Object Lock** (for e-documents/signatures).
  * *Optimization:* **S3 Intelligent-Tiering** to move old scans to cheaper storage automatically.

* **Caching:** **Amazon ElastiCache (Redis)** for session data and hot patient records to reduce DB load.

#### Communication (Edges)

* **Async Messaging (sGW → sPM/sTA):** **Amazon SNS/SQS** for internal decoupling.
* **Real-time Alerts (sTA → iDOC):** **AWS AppSync** (GraphQL subscriptions/WebSockets) or **AWS End User Messaging** (
  replacing Mobile
  Push).
* **Audit Logging (sPM → sAAA):** **Amazon Kinesis Data Streams** for high-volume, "fire-and-forget" audit ingest.

### Development & Maintenance Technologies

AWS tools to support the Software Development Life Cycle (SDLC).

| Phase              | AWS Technology               | Purpose                                                                          |
|--------------------|------------------------------|----------------------------------------------------------------------------------|
| **Source Control** | **AWS CodeCommit**           | Secure, private Git repositories.                                                |
| **CI/CD**          | **AWS CodePipeline**         | Automated testing and deployment to Fargate.                                     |
| **IAC**            | **AWS CDK / CloudFormation** | Infrastructure as Code for environment replication (Selling to other companies). |
| **Monitoring**     | **Amazon CloudWatch**        | Metrics, logs, and alarms for system health.                                     |
| **Observability**  | **AWS X-Ray**                | Distributed tracing to track a request across microservices.                     |
| **Security**       | **AWS Secrets Manager**      | Securely manage DB credentials and API keys.                                     |

### AWS Overview

#### Consolidated Technology Summary

| Official Name                 | Description                                                                               | Documentation Link                                                                                  |
|-------------------------------|-------------------------------------------------------------------------------------------|-----------------------------------------------------------------------------------------------------|
| **AWS AppSync**               | GraphQL/WebSocket service for real-time dashboard updates and alerts.                     | [AppSync Docs](https://docs.aws.amazon.com/appsync/)                                                |
| **AWS CDK**                   | Infrastructure-as-Code framework to define resources using programming languages.         | [CDK Docs](https://docs.aws.amazon.com/cdk/)                                                        |
| **AWS CodeCommit**            | A managed source control service that hosts secure Git-based repositories.                | [CodeCommit Docs](https://docs.aws.amazon.com/codecommit/)                                          |
| **AWS CodePipeline**          | CI/CD service for automating software release workflows.                                  | [CodePipeline Docs](https://docs.aws.amazon.com/codepipeline/)                                      |
| **AWS End User Messaging**    | It provides a unified API for scalable, secure user engagement and transactional alerts.  | [EUM Documentation](https://docs.aws.amazon.com/end-user-messaging/)                                |
| **AWS IoT Core**              | Managed broker for MQTT telemetry ingestion from medical equipment.                       | [IoT Core Docs](https://docs.aws.amazon.com/iot/)                                                   |
| **AWS Lambda**                | Serverless compute for event-driven tasks and light-weight data processing.               | [Lambda Docs](https://docs.aws.amazon.com/lambda/)                                                  |
| **AWS Secrets Manager**       | Securely encrypts, stores, and rotates database credentials, API keys, and other secrets. | [Secrets Manager Docs](https://docs.aws.amazon.com/secretsmanager/)                                 |
| **AWS X-Ray**                 | Distributed tracing to analyze and debug microservice performance.                        | [X-Ray Docs](https://docs.aws.amazon.com/xray/)                                                     |
| **Amazon API Gateway**        | Managed service for creating, publishing, and securing REST APIs.                         | [API Gateway Docs](https://docs.aws.amazon.com/apigateway/)                                         |
| **Amazon Aurora**             | High-performance relational DB for patient records and admin data.                        | [Aurora Docs](https://docs.aws.amazon.com/rds/)                                                     |
| **Amazon CloudFront**         | Content Delivery Network (CDN) to serve the Web Portal with low latency.                  | [CloudFront Docs](https://docs.aws.amazon.com/cloudfront/)                                          |
| **Amazon CloudWatch**         | Monitoring and observability service for logs, metrics, and alarms.                       | [CloudWatch Docs](https://docs.aws.amazon.com/cloudwatch/)                                          |
| **Amazon Comprehend Medical** | NLP service to extract medical information from unstructured clinical text.               | [Comprehend Medical Docs](https://docs.aws.amazon.com/comprehend-medical/)                          |
| **Amazon ECS (Fargate)**      | Serverless container orchestration for core microservices (sPM, sTA, sAAA).               | [ECS Docs](https://docs.aws.amazon.com/ecs/)                                                        |
| **Amazon ElastiCache**        | In-memory caching (Redis) for session management and hot-data access.                     | [ElastiCache Docs](https://docs.aws.amazon.com/elasticache/)                                        |
| **Amazon HealthLake**         | Purpose-built service to store, transform, and analyze healthcare data (FHIR).            | [HealthLake Docs](https://docs.aws.amazon.com/healthlake/)                                          |
| **Amazon Kinesis**            | Real-time data streaming for high-volume audit logs and telemetry.                        | [Kinesis Docs](https://docs.aws.amazon.com/kinesis/)                                                |
| **Amazon Rekognition**        | Deep learning-based computer vision.                                                      | [Rekognition Docs](https://docs.aws.amazon.com/rekognition/)                                        |
| **Amazon S3**                 | Durable object storage for clinical images (DICOM), PDFs, and web hosting.                | [S3 Docs](https://docs.aws.amazon.com/s3/)                                                          |
| **Amazon SNS & SQS**          | Pub/Sub and queuing services for decoupled, asynchronous communication.                   | [SQS Docs](https://docs.aws.amazon.com/sqs/)                                                        |
| **Amazon Timestream**         | Specialized time-series database for high-velocity vitals monitoring.                     | [Timestream Docs](https://docs.aws.amazon.com/timestream/)                                          |
| **Amazon Transcribe Medical** | A HIPAA-eligible speech-to-text service specifically trained on medical vocabulary.       | [Transcribe Medical Docs](https://docs.aws.amazon.com/transcribe/latest/dg/transcribe-medical.html) |

#### SWOT Analysis: AWS AI Technologies

Focusing on clinical value-add and automation.

|            | **Strengths**                                                                                                                        | **Weaknesses**                                                                                             |
|------------|--------------------------------------------------------------------------------------------------------------------------------------|------------------------------------------------------------------------------------------------------------|
| *Internal* | **Amazon HealthLake:** Seamlessly indexes FHIR data.<br><br>**Amazon Comprehend Medical:** Extracts entities from text reports.      | High specialized knowledge required to configure.<br><br>Pricing can scale quickly with high data volumes. |
|            | <center>**Opportunities**</center>                                                                                                   | <center>**Threats**</center>                                                                               |
| *External* | **Amazon Rekognition:** Analyze medical imagery for anomalies.<br><br>**Amazon Transcribe Medical:** Voice-to-text for doctor notes. | Regulatory scrutiny over AI-based diagnosis.<br><br>Potential "Black Box" bias in machine learning models. |

### Deployment Diagram

```mermaid
graph TD
    classDef external fill:#f8a3a3,stroke:#333,stroke-width:1.5px;
    classDef awsService fill:#232F3E,color:#fff,stroke:#232F3E;

%% External Systems
    subgraph External_Systems [External Data Sources]
        Monitors[Bedside Monitors]:::external
        PeerHosp[External Hospital Systems]:::external
    end

%% AWS Cloud
    subgraph AWS_Cloud [AWS Cloud - eu-central-1]
        IGW[Internet Gateway]
        
        subgraph VPC [Virtual Private Cloud - 10.0.0.0/16]
            direction TB
            
            subgraph Public_Subnet [Public Subnet - Ingress]
                direction LR
                ALB[Application Load Balancer]
                IoT[IoT Core]
                NAT[NAT Gateway]
            end

            subgraph Private_Subnet [Private Subnet - Logic]
                direction TB
                Fargate[Fargate Cluster<br/>sAAA, sPM, sTA]
                
                %% PrivateLink Endpoints
                subgraph VPCE [Interface VPC Endpoints / PrivateLink]
                    S3_EP[S3 Endpoint]
                    TS_EP[Timestream Endpoint]
                    KMS_EP[KMS / Secrets Endpoint]
                end

                subgraph Data_Tier [Storage Tier]
                    Aurora[(Aurora DB)]
                    TS[(Timestream DB)]
                    S3[(S3 Bucket)]
                end
            end
        end
    end

%% Connections
    Monitors -->|MQTT/TLS 1.3| IGW
    PeerHosp -->|mTLS/REST| IGW
    IGW --> IoT
    IGW --> ALB
    
    ALB -->|Forward| Fargate
    IoT -->|Ingest| Fargate
    
    %% Internal Routing via PrivateLink
    Fargate --> S3_EP --> S3
    Fargate --> TS_EP --> TS
    Fargate --> KMS_EP
    Fargate -->|ENI| Aurora

    %% NAT for updates/external APIs
    Fargate -.-> NAT -.-> IGW

%% Styling
    style AWS_Cloud fill:#fff,stroke:#FF9900,stroke-width:2px
    style VPC fill:#f9f9f9,stroke:#0073BB,stroke-width:2px
    style Public_Subnet fill:#e1f5fe,stroke:#01579b,stroke-dasharray: 5 5
    style Private_Subnet fill:#f1f8e9,stroke:#2e7d32,stroke-dasharray: 5 5
    style VPCE fill:#fff,stroke:#232f3e,stroke-width:1px
    style Fargate fill:#F58534,color:#fff
    style ALB fill:#8C4FFF,color:#fff
```

This diagram illustrates the **Pulse Patrol Automated Infrastructure**, a cloud-native deployment designed for secure,
high-velocity medical data ingestion and processing within the **AWS eu-central-1** region.

#### Key Architectural Components

* **External Data Sources**: Real-time telemetry is streamed from **Bedside Monitors** via MQTT/TLS 1.3, while
  structured data exchanges with **External Hospital Systems** utilize secure mTLS REST calls.
  * **VPC & Ingress Layer (Public Subnet)**:
  * **VPC (10.0.0.0/16)**: Provides a logically isolated network environment to ensure HIPAA/GDPR compliance.
  * **Internet Gateway (IGW) & NAT Gateway**: The IGW manages external traffic entry, while the NAT Gateway allows
    private resources to securely fetch updates without being exposed to the public internet.
  * **Application Load Balancer (ALB)**: Replaces basic gateway logic to distribute incoming traffic across the
    Fargate cluster for high availability.
  * **IoT Core**: Serves as the managed broker for device connectivity and telemetry routing.


* **Processing Logic (Private Subnet)**:
  * **Fargate Cluster**: Hosts core microservices (**sAAA**, **sPM**, **sTA**) in a private environment, completely
    isolated from direct public access.
  * **VPC Endpoints (PrivateLink)**: Enables the Fargate cluster to communicate with S3, Timestream, and KMS over the
    private AWS backbone, ensuring sensitive PHI never traverses the public internet.


* **Multi-Model Storage Tier**:
  * **Aurora DB**: Stores relational metadata and Electronic Health Records (EHR).
  * **Timestream**: A dedicated time-series database for high-velocity vitals telemetry.
  * **S3 Bucket**: Provides immutable archiving for clinical reports and media recordings, utilizing **Object Lock**
    for legal and regulatory compliance.

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

### API Specifications

#### Patient Transfer API

*Endpoint:* `POST /api/v1/transfers`

*Service:* Patient Management Service (sPM)

*Authentication:* OAuth 2.0 Bearer Token (Admin role required)

*Purpose:* Initiates a secure transfer of patient medical records to an external healthcare provider.

*Request Parameters:*

| Parameter            | Type          | Required | Description                                                                             |
|----------------------|---------------|----------|-----------------------------------------------------------------------------------------|
| `patientId`          | string (UUID) | Yes      | Unique identifier of the patient to transfer                                            |
| `targetProviderId`   | string (UUID) | Yes      | Identifier of the receiving healthcare organization                                     |
| `transferReason`     | string        | Yes      | Clinical justification for transfer (e.g., "Specialist referral", "Patient relocation") |
| `includeTestResults` | boolean       | No       | Include laboratory results (default: true)                                              |
| `includeTelemetry`   | boolean       | No       | Include last 48h of vital signs (default: false)                                        |
| `consentDocumentId`  | string (UUID) | Yes      | Reference to patient consent for data sharing                                           |

*Request Example:*

```json
POST /api/v1/transfers HTTP/1.1
Host: api.pulsepatrol.health
Authorization: Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9...
Content-Type: application/json

{
"patientId": "550e8400-e29b-41d4-a716-446655440000",
"targetProviderId": "7c9e6679-7425-40de-944b-e07fc1f90ae7",
"transferReason": "Patient relocation to new city",
"includeTestResults": true,
"includeTelemetry": false,
"consentDocumentId": "a3bb189e-8bf9-3888-9912-ace4e6543002"
}
```

*Response (Success - 202 Accepted):*

```json        
{
  "transferId": "9f3d7a8b-1c5e-4d2f-8a9b-7e6f5d4c3b2a",
  "status": "processing",
  "estimatedCompletionTime": "2026-02-15T14:35:00Z",
  "trackingUrl": "/api/v1/transfers/9f3d7a8b-1c5e-4d2f-8a9b-7e6f5d4c3b2a/status"
}
```

*Response (Error - 403 Forbidden):*

```json
{
  "error": "INSUFFICIENT_PERMISSIONS",
  "message": "Admin role required to initiate transfers",
  "timestamp": "2026-02-15T14:32:15Z"
}
```

*Response (Error - 409 Conflict):*

```json
{
  "error": "CONSENT_MISSING",
  "message": "Patient consent for data sharing not found or expired",
  "requiredAction": "Obtain patient consent before transfer",
  "timestamp": "2026-02-15T14:32:15Z"
}
```

*Workflow:*

Based on [Sequence 6 (Administrator)](#sequence-6-administrator) diagram

1. Validation Phase:
  - Verify admin authorization token
  - Confirm patient consent is valid and not expired
  - Check target provider exists in peer registry
2. Data Aggregation Phase:
  - Query medical records from Aurora database
  - Fetch laboratory results if requested
  - Retrieve vital signs from Timestream if requested
3. Security Phase:
  - Encrypt PII using AES-256-GCM
  - Generate transfer package with integrity hash (SHA-256)
  - Sign package with hospital's private key
4. Transfer Phase:
  - Send encrypted package to Integration Gateway (sGW)
  - Gateway establishes mTLS tunnel to peer provider
  - Peer provider returns acknowledgment or error
5. Finalization Phase:
  - Mark patient record as "Transferred" in local database
  - Emit audit event to Compliance Service (sAAA)
  - Send completion notification to requesting administrator

*Status Codes:*

| Code | Meaning               | Description                                    |
|------|-----------------------|------------------------------------------------|
| 202  | Accepted              | Transfer initiated successfully                |
| 400  | Bad Request           | Invalid patientId or targetProviderId format   |
| 401  | Unauthorized          | Missing or invalid authentication token        |
| 403  | Forbidden             | User lacks admin privileges                    |
| 404  | Not Found             | Patient or target provider not found           |
| 409  | Conflict              | Missing consent or patient already transferred |
| 500  | Internal Server Error | Database or encryption service failure         |
| 502  | Bad Gateway           | Unable to reach target provider                |

*Rate Limits:*

- 10 transfers per administrator per hour
- 100 transfers per organization per day

*Compliance Notes:*

- All transfers are logged in immutable audit trail (HIPAA requirement)
- Transfer packages are encrypted at rest for 7 years (legal retention)
- Patient can revoke consent retroactively; system will notify receiving provider

## 7. Security Concerns

[//]: # (<<Authorisation, Authentication, Data encryption, Threat modelling diagram>>)

### Data Flow Diagram for  Use Case 6 (Patient Transfer)

```mermaid
graph LR
  subgraph Public_Network ["🌐 Public Internet / External Organizations"]
    Admin>«person»<br/>👤 Administrator]
    Peer(("«software system»<br/>🌐 External Peer<br/>[ePEER]"))
  end

  subgraph AWS_Cloud ["☁️ AWS Cloud - System Boundary"]
    subgraph Public_Subnet ["Public Subnet (DMZ)"]
      WP("«container»<br/>Web Portal<br/>[uiWP]")
      GW("«container»<br/>Integration Gateway<br/>[sGW]")
      IGW["Internet Gateway"]
    end

    subgraph Private_Subnet ["Private Subnet (Logic & Data)"]
      direction TB
      PMS(("«container»<br/>Patient Management<br/>[sPM]"))
      AAA(("«container»<br/>Compliance & Identity<br/>[sAAA]"))

      subgraph Data_Layer ["Storage Tier"]
        DS[("«database»<br/>Patient Records<br/>[pPM]")]
      end
    end
  end

%% Asset Labels (Brown)
  A01{{A01}}:::asset_label --- DS
  A02{{A02}}:::asset_label --- AAA
  A03{{A03}}:::asset_label --- AAA
  A04{{A04}}:::asset_label --- GW
  A05{{A05}}:::asset_label --- PMS

%% Threat Labels (Red)
  TA01{{TA01}}:::threat_label --- IGW
  TA02{{TA02}}:::threat_label --- DS
  TA03{{TA03}}:::threat_label --- AAA
  TA04{{TA04}}:::threat_label --- GW
  TA05{{TA05}}:::threat_label --- PMS

%% Control Labels (Green)
  C01{{C01}}:::control_label --- GW
  C02{{C02}}:::control_label --- DS
  C03{{C03}}:::control_label --- AAA
  C04{{C04}}:::control_label --- WP
  C05{{C05}}:::control_label --- PMS

%% Flow Logic
Admin -- "1. Initiate [HTTPS]" --> IGW
IGW -- "2. API Request [TLS]" --> WP
WP -- "3. Secure Forward" --> PMS
PMS -.->|4. Auth Check - OIDC | AAA
PMS -- "5. Access Personal Health Information" --> DS
PMS -- "6. Transfer Payload" --> GW
GW -- "7. mTLS REST" --> IGW
IGW -- "8. Secure Handshake" --> Peer
PMS -.->|9. Log Transfer| AAA

%% Styling
  classDef boundary fill:none,stroke:#333,stroke-width:2px,stroke-dasharray: 5 5;
  classDef public fill:#e1f5fe,stroke:#01579b;
  classDef private fill:#f1f8e9,stroke:#2e7d32;
  classDef external fill:#f8a3a3,stroke:#333;

%% Label Styles
  classDef asset_label fill:#d7ccc8,stroke:#5d4037,stroke-width:1px,color:#5d4037,font-weight:bold;
  classDef threat_label fill:#DC143C,stroke:#721c24,stroke-width:1px,color:#FFFFFF,font-weight:bold;
  classDef control_label fill:#2E7D32,stroke:#1B5E20,stroke-width:1px,color:#FFFFFF,font-weight:bold;
  
  class AWS_Cloud boundary;
  class Public_Subnet public;
  class Private_Subnet private;
  class Peer external;

%% Link Styling
  linkStyle 0,1,2,3,4 stroke:#5d4037,stroke-width:1px,stroke-dasharray: 1;
  linkStyle 5,6,7,8,9 stroke:#DC143C,stroke-width:1px,stroke-dasharray: 2;
  linkStyle 10,11,12,13,14 stroke:#2E7D32,stroke-width:1px,stroke-dasharray: 4;
```

#### Asset Identification Table

| ID      | Asset Name                | Description                                        | Sensitivity |
|---------|---------------------------|----------------------------------------------------|-------------|
| **A01** | **Patient Records (PHI)** | Highly sensitive medical data stored in `pPM`.     | *Critical*  |
| **A02** | **Identity Tokens**       | OIDC/JWT credentials managed by `sAAA`.            | *High*      |
| **A03** | **Audit Logs**            | Immutable records of access events in `sAAA`.      | *High*      |
| **A04** | **Encryption Keys/Certs** | mTLS and SSL credentials used by `sGW`.            | *High*      |
| **A05** | **Transfer Payloads**     | The transient data packets being moved to `ePEER`. | *Critical*  |

#### Threat Identification Table

| ID       | Threat Category            | Description                                                                                | Targeted Asset(s) |
|----------|----------------------------|--------------------------------------------------------------------------------------------|-------------------|
| **TA01** | **Spoofing**               | Adversary impersonates the `Administrator` or `External Peer` to gain unauthorized access. | A02, A05          |
| **TA02** | **Tampering**              | Unauthorized modification of PHI in transit or at rest in `pPM`.                           | A01, A05          |
| **TA03** | **Repudiation**            | Deletion or modification of audit logs in `sAAA` to hide malicious activity.               | A03               |
| **TA04** | **Info Disclosure**        | Leakage of encryption keys from `sGW` or unencrypted PHI exposure.                         | A01, A04          |
| **TA05** | **Elevation of Privilege** | User bypasses OIDC checks to gain administrative control over `sPM`.                       | A02, A05          |

#### Security Controls & Mitigation Table

| ID      | Control Name              | Description                                                                   | Mitigates  |
|---------|---------------------------|-------------------------------------------------------------------------------|------------|
| **C01** | **mTLS Authentication**   | Mutual TLS for all egress traffic to ensure only verified peers connect.      | TA01, TA04 |
| **C02** | **AES-256 Encryption**    | Encryption at rest and in transit for PHI using AWS KMS managed keys.         | TA02, TA04 |
| **C03** | **Immutable Audit Trail** | Log streaming to sAAA with write-once/read-many (WORM) properties.            | TA03       |
| **C04** | **OIDC & RBAC**           | Identity-based access control via sAAA to prevent vertical movement.          | TA01, TA05 |
| **C05** | **VPC PrivateLink**       | Isolating database traffic within the AWS backbone, bypassing public routing. | TA02, TA04 |

### Compliance with CIA Principles (Confidentiality, Integrity, Availability)

The Pulse Patrol system architecture is engineered to protect **PHI (A01)** by strictly enforcing the CIA triad through
the security controls identified in the threat model:

| Principle           | Implementation Mechanism | Technical Description                                                                                                                                                                                                                                                       |
|---------------------|--------------------------|-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| **Confidentiality** | **C01, C02, C04, C05**   | Access to sensitive records is restricted via **RBAC/OIDC (C04)**. Data is shielded from unauthorized viewing through **AES-256 encryption (C02)** at rest and in transit. Network isolation via **VPC PrivateLink (C05)** ensures PHI never traverses the public internet. |
| **Integrity**       | **C01, C02, C03**        | **mTLS (C01)** ensures that data transferred to the `ePEER` remains unaltered. Any attempt to modify records (TA02) is mitigated by data signing, while the **Immutable Audit Trail (C03)** ensures that the history of the data cannot be falsified or deleted.            |
| **Availability**    | **AWS Infrastructure**   | Availability is maintained through **Multi-AZ deployment** for the `pPM` database and the auto-scaling nature of the Fargate clusters. The **Integration Gateway (sGW)** ensures high-throughput availability for external requests even during peak load.                  |

## 8. COGS

[//]: # (<<Cost estimation model for hardware, services, data storage and transfer for the whole solution>>)

### 8.1 Regions for Comparison

To provide a global perspective on operational costs, the following regions will be used in the AWS Pricing Calculator:

1. **Europe (Frankfurt)** - `eu-central-1`: Primary region (per ARD).
2. **US East (N. Virginia)** - `us-east-1`: Standard baseline for US-based healthcare.

### 8.2 Selected AWS Services for Estimation

The estimation covers the core components of the "Pulse Patrol" infrastructure:

* **Compute:** AWS Fargate (for sPM, sTA, sAAA microservices).
* **Database (Relational):** Amazon Aurora PostgreSQL (for patient records).
* **Database (Time-Series):** Amazon Timestream (for real-time vitals/telemetry).
* **Storage:** Amazon S3 (for clinical reports and immutable logs).
* **Ingress/IoT:** AWS IoT Core (for medical equipment connectivity).
* **Networking:** Amazon CloudFront (for Web Portal and Clinical Dashboard distribution).

### 8.3 Estimation Assumptions

The following assumptions are made for a "Medium-Sized Hospital" deployment:

| Category             | Assumption            | Detail                                                            |
|----------------------|-----------------------|-------------------------------------------------------------------|
| **User Base**        | 5,000 Active Patients | Monthly average of patients under monitoring.                     |
| **Telemetry**        | 100 Bedside Monitors  | Streaming vitals 24/7 at 1 message/second.                        |
| **Compute**          | 6 Fargate Tasks       | 2 tasks per service (sPM, sTA, sAAA) for redundancy.              |
| **Relational Data**  | 500 GB Storage        | Patient history, medical records, and audit logs.                 |
| **Time-Series Data** | 1 TB / Month          | High-resolution physiological data (Heart rate, SpO2).            |
| **Object Storage**   | 2 TB Storage          | Imaging (DICOM), PDFs, and encrypted transfer snapshots.          |
| **Data Transfer**    | 500 GB Outbound       | Traffic from Web Portal and Dashboard to end-users.               |
| **Retention**        | 7 Years               | Long-term archival requirement for legal compliance (S3 Glacier). |

### 8.4 Results

| AWS Service                     | Europe (Frankfurt) | US East (N. Virginia) | 
|---------------------------------|--------------------|-----------------------|
| **AWS Fargate** (Compute)       | $     407.87       | $     354.60          |
| **Amazon Aurora** (Records)     | $     534.99       | $     457.22          |
| **Amazon Timestream** (Vitals)  | $   1,397.22       | $   1,156.32          | 
| **Amazon S3** (Storage/Logs)    | $      52.59       | $      49.26          | 
| **AWS IoT Core** (Connectivity) | $       0.00       | $       0.00          | 
| **Amazon CloudFront** (CDN)     | $     174.08       | $     174.08          | 
| **Other (Networking/KMS)**      | $      79.13       | $      68.88          | 
| **Total Estimated Monthly**     | **$ 2,645.87**     | **$ 2,260.36**        | 

<p>
  <img src="../misc/eu-central-1.png" width="100%">
</p>

<p>
  <img src="../misc/us-east-1.png" width="100%">
</p>