// Implementación del TAD Minecraft.
// Consiste en un Tipo Abstracto de Dato que corre minecraft.
// En el futuro se hara un copia y pega para meterlo al juego mejorado de
// minecraft (zworld).
// Jo profe es que somos de segundo no nos pidas mucho.

template <typename Clave, typename Dato>
struct Minecraft {
    private:
        int blokes;
        int creepers;
        int zombis;
        bool haceFrio;
        struct Minecraft *nodoIzdo;
        struct Minecraft *nodoDcho;

        Clave gatos;
        Dato uwu;
};
