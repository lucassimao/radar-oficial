
## ✅ **MVP Step-by-Step Feature Roadmap**

### 🔹 **Phase 1: Core Infrastructure Setup**

#### 1. **Set up development environment**
- [x] Install and configure **Go 1.24**, keeping Go 1.23 for legacy support
- [x] Set up GitHub repo, environments (local + staging)
- [x] Set up `.env` config with secrets, API keys

#### 2. **Provision and configure infrastructure**
- [ ] Create and configure **PostgreSQL DB** (DO Managed DB or Supabase)
- [x] Create **DigitalOcean AI Knowledge Base**
- [x] Choose embedding model: ✅ `MultiQA MPNet Base Dot v1`

#### 3. **Create Diário ingestion pipeline**
- [ ] Build a Go script to download and normalize Diários Oficiais (PDF → text)
- [ ] Save metadata to PostgreSQL (e.g., date, source, entity)
- [ ] Upload PDF to Object Storage 
- [ ] Trigger update on DO AI KB

---

### 🔹 **Phase 2: WhatsApp Bot + Basic Interaction**

#### 4. **Integrate WhatsApp**
- [ ] Set up webhook to receive/send messages
- [ ] Verify connection and delivery receipt

#### 5. **Create session handling logic**
- [ ] Store basic user session info (phone number, state, entity, current context)
- [ ] Track query counts per user to enforce tier limits

#### 6. **Implement guided conversation flow**
- [ ] Welcome the user with intro message
- [ ] Ask: “De qual estado ou entidade pública você quer buscar informações?”  
- [ ] Store chosen municipality or entity
- [ ] Ask: “Agora me diga o que você gostaria de saber sobre este órgão.”
- [ ] Save both inputs and send a query to DO AI KB

---

### 🔹 **Phase 3: AI Querying and Response**

#### 7. **Send question to DigitalOcean AI KB**
- [ ] Use `/query` API with user prompt + entity filter
- [ ] Retrieve answer and top-matching sources
- [ ] Format response with:
  - Summary
  - Source entity + date
  - Optional snippet from Diário

#### 8. **Send response back to user**
- [ ] Format message for WhatsApp UX
- [ ] Include: “Deseja perguntar mais algo sobre este mesmo órgão ou outro diferente?”

#### 9. **Handle follow-up actions**
- [ ] If same entity: loop to [Step 6]
- [ ] If different entity: reset entity and loop to [Step 5]
- [ ] If user says "não": send thank you + goodbye message

---

### 🔹 **Phase 4: Monetization and User Control**

#### 10. **Implement basic user tiers**
- [ ] Free: 5 questions/month
- [ ] Basic: 5 questions/day
- [ ] Pro: 30 questions/day

#### 11. **Track usage per tier**
- [ ] Use Redis or DB to track message count
- [ ] Block or warn users when limits are hit

#### 12. **Integrate payment via PagSeguro / Pix**
- [ ] Create checkout/payment flow
- [ ] Unlock tier access upon successful payment

---

### 🔹 **Phase 5: Admin and Optimization Tools**

#### 13. **Admin dashboard (optional CLI at MVP)**
- [ ] View user activity
- [ ] View failed queries
- [ ] Manually reset query counts or users

#### 14. **Add basic analytics**
- [ ] Number of users / queries
- [ ] Most searched entities or topics

---

### 💡 Optional Enhancements (Post-MVP)

- [ ] Keyword-based daily alerts via WhatsApp
- [ ] Summarize whole Diário for select entities daily
- [ ] Add support for other states (besides Piauí)
- [ ] Export answer + source to PDF
- [ ] Integrate GPT-4.1 agent for follow-up context & richer conversations

---

Would you like this as a Notion board, Trello template, or markdown checklist to plug into your workspace? I can generate it for you in 1 click.