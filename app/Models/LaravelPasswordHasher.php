<?php

namespace App\Models;

use Illuminate\Hashing\HashManager;

/**
 * Class LaravelPasswordHasher
 *
 * @package App\Models
 */
final class LaravelPasswordHasher implements PasswordHasher
{
    /**
     * @var HashManager
     */
    private HashManager $hashManager;

    /**
     * LaravelPasswordHasher constructor.
     *
     * @param HashManager $hashManager
     */
    public function __construct(HashManager $hashManager)
    {
        $this->hashManager = $hashManager;
    }

    /**
     * @return HashManager
     */
    private function getHashManager(): HashManager
    {
        return $this->hashManager;
    }

    /**
     * @param string $password
     *
     * @return string
     */
    public function hash(string $password): string
    {
        return $this->getHashManager()->make($password);
    }

    /**
     * @param string $password
     * @param string $hashedPassword
     *
     * @return bool
     */
    public function check(string $password, string $hashedPassword): bool
    {
        return $this->getHashManager()->check($password, $hashedPassword);
    }
}
