# Architecture Decision Record - v3.2

> **<span style="font-size: 1.3em">Pulse Patrol</span>**
>
> *Develop a software system for healthcare that collects and manages patient data,
> integrates with medical equipment, provides web access for patients and authorized personnel,
> alerts staff for abnormal values, and supports patient transfers between healthcare providers.*

<!-- TOC -->
* [Architecture Decision Record - v3.2](#architecture-decision-record---v32)
  * [1. Context](#1-context)
    * [Scope](#scope)
    * [Architectural Nodes](#architectural-nodes)
    * [Constraints & Remarks](#constraints--remarks)
  * [2. Alternatives](#2-alternatives)
    * [A. Fully Synchronous Monolith (REST Everywhere)](#a-fully-synchronous-monolith-rest-everywhere)
    * [B. Pure Event-Driven Architecture (Event Sourcing)](#b-pure-event-driven-architecture-event-sourcing)
    * [C. Polling-Based Integration](#c-polling-based-integration)
    * [D. Hybrid Synchronous/Asynchronous with Reactive Push (CHOSEN)](#d-hybrid-synchronousasynchronous-with-reactive-push-chosen)
  * [3. Decision](#3-decision)
    * [Architectural (Communication) Edges](#architectural-communication-edges)
      * [ePP / iADM → uiWP (User Portal Access)](#epp--iadm--uiwp-user-portal-access)
      * [iDOC / iSSM → uiCD (Clinical Monitoring)](#idoc--issm--uicd-clinical-monitoring)
      * [uiWP / uiCD → sAAA (Identity & Auth)](#uiwp--uicd--saaa-identity--auth)
      * [uiWP → sPM (Record Management)](#uiwp--spm-record-management)
      * [uiCD → sTA / sPM (Clinical Data Fetch)](#uicd--sta--spm-clinical-data-fetch)
      * [iEQP → sGW (Telemetry Ingestion)](#ieqp--sgw-telemetry-ingestion)
      * [iLEG ↔ sGW (Legacy Integration)](#ileg--sgw-legacy-integration)
      * [sGW ↔ ePEER (Inter-Provider Transfer)](#sgw--epeer-inter-provider-transfer)
      * [sGW → [PMSq/TASq] → sPM / sTA (Data Normalization)](#sgw--pmsqtasq--spm--sta-data-normalization)
      * [sTA → iDOC / iSSM (Critical Alerting)](#sta--idoc--issm-critical-alerting)
      * [sPM / sTA / sAAA → pDS (Persistence Layer)](#spm--sta--saaa--pds-persistence-layer)
      * [sPM / sTA / sGW → [AAAq] → sAAA (Audit Logging)](#spm--sta--sgw--aaaq--saaa-audit-logging)
  * [4. Consequences](#4-consequences)
<!-- TOC -->

[//]: # (S: <adr>)

## 1. Context

### Scope

[//]: # (<< What is the problem we are trying to solve >>)

The scope of the document is to serve as a reference for the decisions related to inter-service communication
in the Pulse Patrol software solution.

It will:

- reproduce the main architectural nodes
- iterate over communication technologies
- assign edges (communication technologies) between any two nodes that are supposed to communicate

### Architectural Nodes

**Environmental**

- Human Actors
    - *External*
        - **ePP - Patients** - views personal medical history, test results, and treatment progress via the web portal
    - *Internal*
        - **iDOC - Doctor** - accesses patient data within their hospital and receives critical physiological alerts
        - **iSSM - Support Staff Member** - nurses/assistants who receive real-time alerts for abnormal patient
          monitoring values
        - **iADM - Administrator** - manages records, oversees data integrity, and initiates inter-company patient
          transfers
- Technical Systems
    - *External*
        - **ePEER - External Healthcare Companies Peer** - systems belonging to other providers that receive or send
          patient data during a transfer
    - *Internal*
        - **iEQP - Medical Equipment** - IoT devices and monitoring hardware (e.g., bedside monitors, ventilators) that
          stream real-time telemetry
        - **iLEG - Legacy Hospital Systems** - existing legacy databases or EMRs where admission forms and historical
          medical records may reside

**Containers**

- **uiWP - Web Portal**: Interface for Patients to view records and for Administrators to manage data.
- **uiCD - Clinical Dashboard**: Specialized interface for Doctors and Support Staff to monitor live telemetry and
  patient data.
- **sGW - Integration Gateway**: Handles communication with Legacy Systems, Medical Equipment, and Peer Healthcare
  Companies.
- **sPM - Patient Management Services**: Core logic for medical records, admission forms, and inter-company transfers.
- **sTA - Telemetry & Alerting Services**: Processes real-time data from medical equipment and triggers notifications
  for abnormal values.
- **sAAA - Compliance & Identity Services**: Handles authentication, authorization, audit, ...
- **pDS - Data Storage**: Centralized repository for structured medical records and time-series telemetry data.

### Constraints & Remarks

To ensure the technical feasibility and regulatory compliance of the Pulse Patrol architecture, the following
constraints and operational principles are applied:

- _Database Schema Isolation_: To simplify the documentation of persistence, compute-to-storage connections are
  considered
  exclusive to the service's own schema. sPM and sTA do not share database tables; all cross-domain data exchange must
  occur via the defined communication edges.

- _Zero-Trust Identity_: Every synchronous request across the system (internal or external) must be accompanied by a
  valid
  identity token (OIDC/OAuth2). No service trusts another based solely on network location.

- _Telemetry Integrity (QoS)_: Given the critical nature of medical data, MQTT communication from iEQP must utilize QoS
  Level 2 (Exactly Once) for alerts and QoS Level 0 (At Most Once) for high-frequency wave-form data to balance
  reliability with network bandwidth.

- _Regulatory Data Locality_: All patient-identifiable information (PII) handled by sPM and sGW must be encrypted at
  rest
  and in transit. During ePEER transfers, data must comply with regional healthcare data residency laws (e.g., HIPAA,
  GDPR).

- _Non-Blocking Audit Path_: The audit logging mechanism must be "fire-and-forget" from the perspective of the primary
  services (sPM, sTA, sGW). A failure in the sCI audit queue must not halt clinical operations, though it must trigger
  an
  immediate administrative system alert.

- _Legacy Protocol Translation_: The sGW is responsible for all protocol translation. Internal services will only
  communicate using modern standards (JSON/REST/AsyncAPI), shielding the core logic from legacy HL7 v2 or flat-file
  formats.

In order to simplify the communication, for compute to persistence connections only communication with its own schema
will be considered.
For example if a link between **sAAA** and **pDS** is mentioned that exclusively implies a connection between the
Compliance & Identity Services and the corresponding schema.

## 2. Alternatives

[//]: # (<< What are the alternatives considered >>)

Before settling on the hybrid synchronous/asynchronous architecture detailed in Section 3, the following three
alternative communication strategies were evaluated:

### A. Fully Synchronous Monolith (REST Everywhere)

All components communicate via direct HTTP/REST calls.

**Pros**: Simpler to develop and debug; strong immediate consistency.

**Cons**: High Risk. If the Audit Service (sAAA) or Data Storage (pDS) lags, the entire telemetry stream from iEQP
stalls.
In a clinical setting, this latency could delay life-saving alerts. It lacks the buffering needed for high-frequency
medical data.

### B. Pure Event-Driven Architecture (Event Sourcing)

Every action, from a login to a heart rate change, is an event stored in a central log (e.g., Kafka).

**Pros**: Perfect audit trail by default; highly decoupled services.

**Cons**: Overly Complex. Doctors and Patients expect immediate feedback when updating records (Read-after-Write
consistency). Eventual consistency in a medical record environment could lead to "stale" data (e.g., showing an old
medication list right after an update), which is a significant safety risk.

### C. Polling-Based Integration

Instead of WebSockets or Push Notifications, the UI and Gateway "poll" the services for updates at regular intervals (
e.g., every 5 seconds).

**Pros**: Easier to implement with standard legacy load balancers.

**Cons**: Inefficient. High battery drain on mobile devices for iSSM/iDOC and unacceptable delays for critical alerts. A
5-second polling delay is too long for a cardiac arrest notification.

### D. Hybrid Synchronous/Asynchronous with Reactive Push (CHOSEN)

Combines synchronous REST/gRPC for user operations, asynchronous message queues
for high-volume data, and reactive WebSocket push for real-time alerts.

**Pros**: Balances immediate consistency for medical records with high-throughput
for telemetry. Real-time push prevents alert delays.

**Cons**: More complex infrastructure requiring multiple technologies.

## 3. Decision

[//]: # (<< What alternative was chosen and why >>)

The **Hybrid Architecture** (D) was chosen to balance the immediate consistency required for medical records with the
high-throughput, non-blocking requirements of medical telemetry and audit logging.

- **Synchronous (REST/HTTPS)** is used for user-facing actions where immediate confirmation is required.

- **Asynchronous (Message Bus/Queue)** is used for telemetry and audit logs to ensure that spikes in data volume do not
  crash
  the user interface and that logging never slows down clinical workflows.

- **Reactive (WebSockets/Push)** ensures that "Critical Alerts" are pushed to staff immediately rather than waiting for
  a
  manual screen refresh.

### Architectural (Communication) Edges

**Remark:** _the <u>source of the arrow</u> indicates the part that <u>initiates the communication</u>._

#### ePP / iADM → uiWP (User Portal Access)

_Type_: Synchronous, Request/Response, (HTTPS/TLS (REST))

_Detail_: Patients and Admins initiate sessions via browser. High-security measures (MFA) are enforced.

#### iDOC / iSSM → uiCD (Clinical Monitoring)

_Type_: Sync & Reactive, Bidirectional, (HTTPS + WebSockets (WSS))

_Detail_: Clinical staff initiate the connection; uiCD maintains a persistent socket to receive live telemetry updates
pushed from the backend.

#### uiWP / uiCD → sAAA (Identity & Auth)

_Type_: Synchronous, Request/Response, (OIDC / OAuth2 (over HTTPS))

_Detail_: UI components initiate authentication and authorization requests to verify user permissions.

#### uiWP → sPM (Record Management)

_Type_: Synchronous, Request/Response, (REST (JSON))

_Detail_: Web Portal initiates requests for medical record views, admission form submissions, and transfer triggers.

#### uiCD → sTA / sPM (Clinical Data Fetch)

_Type_: Synchronous , Request/Response, (REST or gRPC)

_Detail_: Dashboard initiates data fetches for patient history (sPM) and real-time monitoring configurations (sTA).

#### iEQP → sGW (Telemetry Ingestion)

_Type_: Asynchronous, Stream, (MQTT (QoS 0/2))

_Detail_: Medical equipment initiates the data flow, streaming raw telemetry packets to the gateway.

#### iLEG ↔ sGW (Legacy Integration)

_Type_: Asynchronous, File/Batch/Stream, (HL7 FHIR / MLLP)

_Detail_: sGW initiates pulls for historical data, while iLEG may push updates regarding new admissions.

#### sGW ↔ ePEER (Inter-Provider Transfer)

_Type_: Synchronous, Request/Response, (mTLS REST API)

_Detail_: Handshake-based transfer where either peer can initiate a record transfer request.

#### sGW → [PMSq/TASq] → sPM / sTA (Data Normalization)

_Type_: Asynchronous, Queue, (Message Broker specific)

_Detail_: The gateway pushes normalized data to intermediate queues (PMSq for patient management, TASq for telemetry).
Services consume from their respective queues for processing and storage.

#### sTA → iDOC / iSSM (Critical Alerting)

_Type_:  Asynchronous, Push Notification, (WebSockets / Mobile Push Service)

_Detail_: sTA initiates the flow upon detecting abnormal vitals, alerting staff immediately.

#### sPM / sTA / sAAA → pDS (Persistence Layer)

_Type_: Synchronous, Direct Connection, (DB Driver / gRPC)

_Detail_: Services initiate read/write operations to their respective database schemas.

#### sPM / sTA / sGW → [AAAq] → sAAA (Audit Logging)

_Type_: Asynchronous, Queue, (Message Broker specific)

_Detail_: Every time a patient data access occurs (e.g., a Doctor viewing a record, a new lab report coming in), the
handling service emits a fire-and-forget event to the Audit Queue (AAAq) containing the Actor ID, Patient ID, Timestamp,
and Resource Accessed. The sAAA service consumes events from this queue.

## 4. Consequences

[//]: # (<< What are the implications of our decision? What components/systems are impacted and how? >>)

The adoption of a hybrid synchronous/asynchronous communication model significantly impacts the system's operational
profile, reliability, and development lifecycle.

**Positive (Benefits)**

- System Responsiveness: By offloading Audit Logging and Telemetry Processing to asynchronous workers, the UI remains
  highly responsive. Doctors do not experience "spinner fatigue" while waiting for background compliance tasks to
  complete.

- Scalability & Resilience: The use of a Message Bus between sGW and core services decouples high-volume IoT data from
  business logic. If sPM undergoes maintenance, sGW can continue to buffer incoming telemetry, preventing data loss.

- Regulatory Compliance: Asynchronous audit events ensure a comprehensive "paper trail" of patient data access without
  compromising the performance of clinical systems.

- Operational Safety: Real-time Push Notifications (WebSockets/Mobile Push) minimize the "Time-to-Alert," directly
  improving patient outcomes during critical physiological events.

**Negative (Trade-offs)**

- Increased Infrastructure Complexity: The architecture requires managing a diverse tech stack, including an MQTT
  Broker (
  for devices), a Message Broker (for internal events), and a Push Gateway (for mobile alerts).

- Eventual Consistency in Analytics: While medical records in sPM are strictly consistent, the telemetry visualizations
  and audit logs are "eventually consistent." There may be a sub-second delay between a data event occurring, and it's
  appearing in the audit trail.

- Monitoring Burden: Distributed systems are harder to monitor. We must implement Distributed Tracing (e.g.,
  OpenTelemetry) to track a single patient's data journey from the iEQP through the sGW to the sTA.

**Neutral (Operational Requirements)**

- Specialized Storage: The decision to use a Time-Series Driver for pDS requires the DevOps team to maintain two
  distinct
  database technologies (Relational and Time-Series) to handle the different data velocities.

- Strict Identity Management: Since the system is highly distributed, sAAA becomes a single point of failure (SPOF) for
  security. High availability (HA) for the identity service is mandatory.

[//]: # (S: </adr>)