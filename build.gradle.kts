plugins {
    kotlin("jvm") version "1.9.21"
}

group = "org.aria.danesh"
version = "1.0-SNAPSHOT"

repositories {
    mavenCentral()
}

dependencies {
    implementation(project(mapOf("path" to ":LogManagingKotlinLib")))
    testImplementation("org.jetbrains.kotlin:kotlin-test")
    implementation("org.jetbrains.kotlinx:kotlinx-coroutines-core:1.7.3")
    implementation("org.json:json:20231013")
}

tasks.test {
    useJUnitPlatform()
}
kotlin {
    jvmToolchain(17) // Or your desired JVM version
}

tasks.jar {
    manifest {
        attributes["Main-Class"] = "org.aria.danesh.logmanagingkotlinlib.LogManagingKotlinLibKt" // If you have a main function. if not, remove this line.
    }
    from(configurations.runtimeClasspath.get().map { if (it.isDirectory) it else zipTree(it) })
}