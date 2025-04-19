import React from 'react';

const Header = () => {
  return (
    <header className="fixed w-full bg-white bg-opacity-90 backdrop-blur-sm z-50 shadow-sm">
      <div className="max-w-screen-lg mx-auto px-4 py-4 flex justify-between items-center">
        <div className="flex items-center">
          <h1 className="text-2xl font-bold text-blue-600">Radar Oficial</h1>
        </div>
        <nav className="hidden md:flex space-x-8">
          <a href="#features" className="text-gray-700 hover:text-blue-600 transition-colors duration-200">
            Funcionalidades
          </a>
          <a href="#showcase" className="text-gray-700 hover:text-blue-600 transition-colors duration-200">
            Demo
          </a>
          <a href="#pricing" className="text-gray-700 hover:text-blue-600 transition-colors duration-200">
            Preços
          </a>
          <a href="mailto:contato@radaroficial.app" className="text-gray-700 hover:text-blue-600 transition-colors duration-200">
            Contato
          </a>
        </nav>
        <div className="md:hidden">
          <button className="text-gray-700 focus:outline-none">
            <svg xmlns="http://www.w3.org/2000/svg" className="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 6h16M4 12h16M4 18h16" />
            </svg>
          </button>
        </div>
      </div>
    </header>
  );
};

export default Header;
