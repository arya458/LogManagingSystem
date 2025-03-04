plugins {
    id("java-library")
    id("org.jetbrains.kotlin.jvm")
}

java {
    sourceCompatibility = JavaVersion.VERSION_17 // Or a higher version if needed
    targetCompatibility = JavaVersion.VERSION_17
}
repositories {
    mavenCentral()
}

dependencies {
    implementation("org.jetbrains.kotlin:kotlin-stdlib-jdk8") // Or jdk7 if needed
    implementation("org.json:json:20231013") // Add JSON dependency
    implementation("org.jetbrains.kotlinx:kotlinx-coroutines-core:1.7.3") // Coroutines core
    implementation("org.jetbrains.kotlinx:kotlinx-coroutines-jdk8:1.7.3") // Coroutines jdk8
}

tasks.test {
    useJUnitPlatform()
}
kotlin {
    jvmToolchain(17)
}