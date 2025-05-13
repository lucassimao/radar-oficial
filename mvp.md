## ✅ **Updated MVP Feature Roadmap**

### 🔹 **Phase 1: Core Infrastructure Setup**

#### 1. **Set up development environment**
- [x] Install and configure **Go 1.24**, keeping Go 1.23 for legacy support  
- [x] Set up GitHub repo, environments (local + staging)  
- [x] Set up `.env` config with secrets, API keys  

#### 2. **Provision and configure infrastructure**
- [x] Create and configure **PostgreSQL DB**  
- [x] Create **DigitalOcean AI Knowledge Base**  
- [x] Choose embedding model: ✅ `MultiQA MPNet Base Dot v1`  

#### 3. **Create Diário ingestion pipeline**
- [x] Build a Go script to download and normalize Diários Oficiais (PDF → text)  
- [x] Save metadata to PostgreSQL (e.g., date, source, entity)  
- [x] Upload PDF to Object Storage  
- [x] Trigger update on DO AI KB  

---

### 🔹 **Phase 2: Web Chatbot MVP + Query Management**

#### 4. **Create basic chatbot web interface**
- [x] Build a lightweight web UI with Next.js (chat style like ChatGPT)
- [ ] Welcome message + entity selection flow
        - if no previous sessions, send welcome msg
            - select entity
            - save entity in the local state
            - send msg + target entity
        - Subscribe

        - find a way to allow change the selected diario
        - if users cross the free usage limits, then display a msg asking for payment

- [x] Allow user to input natural language queries
- [ ] Show response from DO AI KB + source metadata
- [ ] Refactor river workers
- [ ] Fetch diarios of all states
- [ ] Display diarios in the website

#### 5. **Handle user sessions and query usage**
- [ ] Store user ID or session fingerprint
- [ ] Track number of queries per user
- [ ] Enforce free tier limit (5/month)

---

### 🔹 **Phase 3: AI Querying and Response**

#### 6. **Send question to DigitalOcean AI KB**
- [ ] Use `/query` API with user prompt + entity filter  
- [ ] Retrieve answer and top-matching sources  
- [ ] Allow PDF downloads

#### 7. **Send response back to frontend**
- [ ] Format answer for chatbot UX
- [ ] Ask user: "Deseja continuar com este mesmo órgão ou outro?"

#### 8. **Handle follow-up flow**
- [ ] If same entity: loop
- [ ] If new entity: reset and start again
- [ ] If "não": close session with thank-you

---

### 🔹 **Phase 4: WhatsApp Integration**

#### 9. **Integrate WhatsApp API**
- [x] Set up webhook for incoming messages
- [x] Send & receive messages
- [ ] Setup new Whatsapp Account
- [ ] Mirror the chatbot flow in WhatsApp

#### 10. **Connect to query tracking logic**
- [ ] Reuse query counting mechanism for WhatsApp users
- [ ] Enforce plan limits and usage tracking

---

### 🔹 **Phase 5: Monetization**

#### 11. **Implement user tiers**
- [ ] Free: 5 queries/month  
- [ ] Basic: 5 queries/day  
- [ ] Pro: 30 queries/day  

#### 12. **Integrate Pix / PagSeguro payments**
- [ ] Payment via Pix (e.g. Gerencianet, Asaas, Mercado Pago)  
- [ ] Webhook to confirm payment and unlock tier access  
- [ ] Store and manage plan assignments per user

---

### 🔹 **Phase 6: Admin & Analytics**

#### 13. **Admin tools (CLI or web)**
- [ ] View usage stats  
- [ ] View failed queries  
- [ ] Reset users / limits
- [ ] Setup River jobs dashboard


#### 14. **Basic analytics**
- [ ] Total queries / users  
- [ ] Most searched entities or terms

---

### 💡 **Optional Post-MVP Enhancements**

- [ ] Keyword-based alerts via WhatsApp  
- [ ] Diário summarization per entity  
- [ ] Support more states beyond Piauí  
- [ ] GPT-4.1 for richer, context-aware chat  
- [ ] Export chat + Diário to PDF  
