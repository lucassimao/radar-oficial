import React from "react";

import Image from "next/image";

const Hero = () => {
  return (
    <section className="min-h-screen flex items-center pt-20 pb-16 px-4">
      <div className="max-w-screen-lg mx-auto">
        <div className="flex flex-col md:flex-row items-center">
          <div className="md:w-1/2 md:pr-12">
            <div className="inline-block mb-6">
              <span className="bg-blue-100 text-blue-800 text-sm font-medium px-4 py-2 rounded-full flex items-center">
                🚀 Lançamento previsto para Junho de 2025
              </span>
            </div>
            <h1 className="text-4xl md:text-5xl lg:text-6xl font-bold leading-tight mb-6 text-gray-900">
              Acompanhe os Diários Oficiais com inteligência
            </h1>
            <p className="text-xl text-gray-600 mb-8 leading-relaxed">
              Radar Oficial reúne os Diários Oficiais de todo o Brasil em um só
              lugar — com busca inteligente, notificações por e-mail baseadas em
              palavras-chave e acesso via uma interface web conversacional, no
              estilo de um chatbot como o ChatGPT.
            </p>
            <a
              href="#features"
              className="bg-blue-600 hover:bg-blue-700 text-white font-medium py-3 px-8 rounded-lg shadow-lg hover:shadow-xl transition-all duration-200 inline-block"
            >
              Ver funcionalidades
            </a>
          </div>
          <div className="md:w-1/2 mt-12 md:mt-0 hidden md:block">
            <div className="bg-gradient-to-br from-blue-500 to-blue-700 rounded-2xl h-96 flex items-center justify-center shadow-2xl hover:scale-105 transition-all duration-300 overflow-hidden">
              <Image
                src="/images/hero2.png" // 👈 salve a imagem em public/hero-illustration.png
                alt="Ilustração representando o Radar Oficial"
                width={400}
                height={400}
                className="object-contain"
                priority
              />
            </div>
          </div>
        </div>
      </div>
    </section>
  );
};

export default Hero;
