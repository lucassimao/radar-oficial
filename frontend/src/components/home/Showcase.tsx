import React from 'react';

const Showcase = () => {
  return (
    <section id="showcase" className="py-20 bg-gray-50">
      <div className="max-w-screen-lg mx-auto px-4">
        <div className="max-w-4xl mx-auto">
          <div className="bg-white p-8 rounded-2xl shadow-lg overflow-hidden mb-10">
            <div className="bg-gray-100 rounded-xl p-6 flex flex-col md:flex-row items-center gap-6">
              <div className="bg-green-500 rounded-full w-14 h-14 flex items-center justify-center flex-shrink-0">
                <svg xmlns="http://www.w3.org/2000/svg" className="h-8 w-8 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 10h.01M12 10h.01M16 10h.01M9 16H5a2 2 0 01-2-2V6a2 2 0 012-2h14a2 2 0 012 2v8a2 2 0 01-2 2h-5l-5 5v-5z" />
                </svg>
              </div>
              <div className="space-y-3">
                <div className="bg-blue-100 text-blue-800 rounded-full px-4 py-1 text-sm inline-block">Usuário</div>
                <p className="bg-blue-50 p-3 rounded-lg text-gray-700">Como faço para saber se meu nome apareceu no Diário Oficial de Teresina?</p>
                
                <div className="bg-green-100 text-green-800 rounded-full px-4 py-1 text-sm inline-block">Radar Oficial</div>
                <p className="bg-green-50 p-3 rounded-lg text-gray-700">
                  Encontrei 3 menções ao seu nome nos últimos 15 dias no Diário Oficial de Teresina. A mais recente foi na edição de 10/05/2025, página 45, referente ao processo administrativo 2025.004.3721.
                </p>
              </div>
            </div>
          </div>
          
          <div className="text-center">
            <h2 className="text-3xl font-bold mb-4">Você pergunta. Radar Oficial responde.</h2>
            <p className="text-xl text-gray-600">Simples, direto e sem complicação.</p>
          </div>
        </div>
      </div>
    </section>
  );
};

export default Showcase;