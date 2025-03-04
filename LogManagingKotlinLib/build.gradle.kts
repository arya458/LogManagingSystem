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
tasks.build {
    exec {
        commandLine ("${rootDir}\\build.bat", "windows", "amd64", "LogManagingApi-windows-64.exe")
    }
//    exec {
//        commandLine ("${rootDir}\\build.bat", "windows", "386", "LogManagingApi-windows-32.exe")
//    }
//    exec {
//        commandLine ("${rootDir}\\build.bat", "linux", "amd64", "LogManagingApi_linux_amd64")
//    }
//    exec {
//        commandLine ("${rootDir}\\build.bat", "linux", "386", "LogManagingApi_linux_386")
//    }
//    exec {
//        commandLine ("${rootDir}\\build.bat", "linux", "arm64", "LogManagingApi_linux_arm64")
//    }
//    exec {
//        commandLine ("${rootDir}\\build.bat", "linux", "arm", "LogManagingApi_linux_arm")
//    }
//    exec {
//        commandLine ("${rootDir}\\build.bat", "darwin", "amd64", "LogManagingApi_darwin_amd64")
//    }
//    exec {
//        commandLine("${rootDir}\\build.bat", "darwin", "arm64", "LogManagingApi_darwin_arm64")
//    }
}