import React from "react";

const Footer = () => {
  return (
    <footer className="bg-gray-900 text-white py-10">
      <div className="max-w-screen-lg mx-auto px-4">
        <div className="flex flex-col md:flex-row justify-between items-center md:items-start">
          <div className="mb-8 md:mb-0 text-center md:text-left">
            <h2 className="text-2xl font-bold mb-2">Radar Oficial</h2>
          </div>

          <div className="flex flex-col md:flex-row gap-4 md:gap-8 items-center mb-8 md:mb-0">
            <a
              href="#"
              className="text-gray-400 hover:text-white transition-colors duration-200"
            >
              Termos de Uso
            </a>
            <a
              href="#"
              className="text-gray-400 hover:text-white transition-colors duration-200"
            >
              Política de Privacidade
            </a>
            <a
              href="mailto:contato@radaroficial.app"
              className="text-gray-400 hover:text-white transition-colors duration-200"
            >
              Contato
            </a>
          </div>
        </div>

        <div className="mt-8 pt-8 border-t border-gray-800 text-center md:text-left text-gray-400">
          <p className="mb-2">
            © {new Date().getFullYear()} Radar Oficial. Todos os direitos
            reservados.
          </p>
          <p className="max-w-3xl text-sm">
            Radar Oficial é um produto independente para facilitar o acesso à
            informação pública de todo o Brasil.
          </p>
          <p className="mt-2 text-sm">
            🚀 Em desenvolvimento - Lançamento previsto para Junho de 2025
          </p>
        </div>
      </div>
    </footer>
  );
};

export default Footer;
